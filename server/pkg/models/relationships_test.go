package models

import (
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&User{}, &Organisation{}, &Cars{}, &Rides{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create valid user", func(t *testing.T) {
		user := User{
			Fullname: "John Doe",
			Email:    "john@example.com",
			Phone:    "1234567890",
			Password: "password123",
			Rank:     "user",
		}

		result := db.Create(&user)
		if result.Error != nil {
			t.Errorf("Failed to create user: %v", result.Error)
		}

		if user.ID == 0 {
			t.Error("User ID should not be 0 after creation")
		}
	})

	t.Run("Create user with invalid rank", func(t *testing.T) {
		user := User{
			Fullname: "Jane Doe",
			Email:    "jane@example.com",
			Phone:    "0987654321",
			Password: "password456",
			Rank:     "invalid_rank",
		}

		result := db.Create(&user)
		if result.Error == nil {
			t.Error("Expected an error when creating user with invalid rank")
		} else if result.Error.Error() != "invalid rank value" {
			t.Errorf("Expected 'invalid rank value' error, got: %v", result.Error)
		}
	})

	t.Run("Create user with duplicate email", func(t *testing.T) {
		user1 := User{
			Fullname: "User One",
			Email:    "duplicate@example.com",
			Phone:    "1111111111",
			Password: "password1",
			Rank:     "user",
		}

		user2 := User{
			Fullname: "User Two",
			Email:    "duplicate@example.com",
			Phone:    "2222222222",
			Password: "password2",
			Rank:     "user",
		}

		db.Create(&user1)
		result := db.Create(&user2)
		if result.Error == nil {
			t.Error("Expected an error when creating user with duplicate email")
		}
	})
}

func TestOrganisationModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create organisation with owner", func(t *testing.T) {
		owner := User{
			Fullname: "Owner",
			Email:    "owner@example.com",
			Phone:    "1234567890",
			Password: "ownerpass",
			Rank:     "owner",
		}

		org := Organisation{
			Name:       "Test Org",
			ProfilePic: "profile.jpg",
			Owner:      owner,
		}

		result := db.Create(&org)
		if result.Error != nil {
			t.Errorf("Failed to create organisation: %v", result.Error)
		}

		if org.ID == 0 {
			t.Error("Organisation ID should not be 0 after creation")
		}

		if org.OwnerID == 0 {
			t.Error("Organisation OwnerID should not be 0 after creation")
		}
	})

	t.Run("Add participants to organisation", func(t *testing.T) {
		org := Organisation{
			Name:       "Org with Participants",
			ProfilePic: "org_profile.jpg",
		}

		db.Create(&org)

		participants := []User{
			{Fullname: "User1", Email: "user1@example.com", Phone: "1111111111", Password: "pass1", Rank: "user"},
			{Fullname: "User2", Email: "user2@example.com", Phone: "2222222222", Password: "pass2", Rank: "editor"},
		}

		for i := range participants {
			participants[i].OrganisationID = &org.ID
		}

		result := db.Create(&participants)
		if result.Error != nil {
			t.Errorf("Failed to create participants: %v", result.Error)
		}

		var fetchedOrg Organisation
		db.Preload("Participants").First(&fetchedOrg, org.ID)

		if len(fetchedOrg.Participants) != 2 {
			t.Errorf("Expected 2 participants, got %d", len(fetchedOrg.Participants))
		}
	})
}

func TestUserOrganisationRelationship(t *testing.T) {
	db := setupTestDB(t)

	t.Run("User owns organisation", func(t *testing.T) {
		owner := User{
			Fullname: "Owner User",
			Email:    "owner@test.com",
			Phone:    "9876543210",
			Password: "ownerpass",
			Rank:     "owner",
		}

		// Create the owner first
		if err := db.Create(&owner).Error; err != nil {
			t.Fatalf("Failed to create owner: %v", err)
		}

		org := Organisation{
			Name:       "Owned Org",
			ProfilePic: "owned_org.jpg",
			OwnerID:    owner.ID,
		}

		// Create the organisation
		if err := db.Create(&org).Error; err != nil {
			t.Fatalf("Failed to create organisation: %v", err)
		}

		// Update the owner with the OwnedOrganisationID
		owner.OwnedOrganisationID = &org.ID
		if err := db.Save(&owner).Error; err != nil {
			t.Fatalf("Failed to update owner: %v", err)
		}

		var fetchedUser User
		err := db.Preload("OwnedOrganisation").First(&fetchedUser, owner.ID).Error
		if err != nil {
			t.Fatalf("Failed to fetch user: %v", err)
		}

		if fetchedUser.OwnedOrganisation == nil {
			t.Error("User should have an owned organisation")
		} else if fetchedUser.OwnedOrganisation.ID != org.ID {
			t.Error("User's owned organisation ID doesn't match the created organisation")
		}
	})

	t.Run("User belongs to organisation", func(t *testing.T) {
		org := Organisation{
			Name:       "Parent Org",
			ProfilePic: "parent_org.jpg",
		}

		db.Create(&org)

		user := User{
			Fullname:       "Member User",
			Email:          "member@test.com",
			Phone:          "5555555555",
			Password:       "memberpass",
			Rank:           "user",
			OrganisationID: &org.ID,
		}

		db.Create(&user)

		var fetchedUser User
		db.Preload("Organisation").First(&fetchedUser, user.ID)

		if fetchedUser.Organisation == nil {
			t.Error("User should belong to an organisation")
		}

		if fetchedUser.Organisation.ID != org.ID {
			t.Error("User's organisation ID doesn't match the created organisation")
		}
	})
}

func TestCarModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create car belonging to a user", func(t *testing.T) {
		user := User{
			Fullname: "Car Owner",
			Email:    "carowner@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}

		db.Create(&user)

		car := Cars{
			Name:         "Toyota Carola",
			Type:         "pkw",
			LicensePlate: "ABC123",
			UserID:       &user.ID,
		}

		result := db.Create(&car)
		if result.Error != nil {
			t.Errorf("Failed to create car: %v", result.Error)
		}

		var fetchedUser User
		db.Preload("Cars").First(&fetchedUser, user.ID)

		if len(fetchedUser.Cars) != 1 {
			t.Errorf("Expected user to have 1 car, got %d", len(fetchedUser.Cars))
		}
	})

	t.Run("Create car belonging to an organisation", func(t *testing.T) {
		org := Organisation{
			Name:       "Car Org",
			ProfilePic: "car_org.jpg",
		}

		db.Create(&org)

		car := Cars{
			Name:           "Honda Civic",
			Type:           "pkw",
			LicensePlate:   "XYZ789",
			OrganisationID: &org.ID,
		}

		result := db.Create(&car)
		if result.Error != nil {
			t.Errorf("Failed to create car: %v", result.Error)
		}

		var fetchedOrg Organisation
		db.Preload("Cars").First(&fetchedOrg, org.ID)

		if len(fetchedOrg.Cars) != 1 {
			t.Errorf("Expected organisation to have 1 car, got %d", len(fetchedOrg.Cars))
		}
	})

	t.Run("Car cannot belong to both user and organisation", func(t *testing.T) {
		user := User{
			Fullname: "Another User",
			Email:    "another@example.com",
			Phone:    "9876543210",
			Password: "password456",
			Rank:     "user",
		}

		org := Organisation{
			Name:       "Another Org",
			ProfilePic: "another_org.jpg",
		}

		db.Create(&user)
		db.Create(&org)

		car := Cars{
			Name:           "Ford Mustang",
			Type:           "pkw",
			LicensePlate:   "DEF456",
			UserID:         &user.ID,
			OrganisationID: &org.ID,
		}

		result := db.Create(&car)
		if result.Error == nil {
			t.Error("Expected an error when creating car belonging to both user and organisation")
		} else if !strings.Contains(result.Error.Error(), "car cannot belong to both a user and an organization") {
			t.Errorf("Unexpected error: %v", result.Error)
		}
	})
}

func TestRideModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create ride associated with a car", func(t *testing.T) {
		// Create a user
		user := User{
			Fullname: "Ride User",
			Email:    "rideuser@example.com",
			Phone:    "1234567890",
			Password: "password123",
			Rank:     "user",
		}
		db.Create(&user)

		// Create a car
		car := Cars{
			Name:         "Tesla Model 3",
			Type:         "pkw",
			LicensePlate: "EV123",
			UserID:       &user.ID,
		}
		db.Create(&car)

		// Create a ride
		description := "A test ride"
		rideFrom := "Home"
		rideTo := "Work"
		distance := "30km"
		ride := Rides{
			Title:       "Morning Commute",
			Description: &description,
			RideFrom:    &rideFrom,
			RideTo:      &rideTo,
			Stops:       []string{"Coffee Shop", "Gas Station"},
			Distance:    &distance,
			BeginTime:   "2024-08-02T08:00:00Z",
			EndTime:     "2024-08-02T09:00:00Z",
			Category:    "Work",
			CarID:       car.ID,
		}

		result := db.Create(&ride)
		if result.Error != nil {
			t.Errorf("Failed to create ride: %v", result.Error)
		}

		// Fetch the ride and check its associations
		var fetchedRide Rides
		db.Preload("Car").First(&fetchedRide, ride.ID)

		if fetchedRide.Car.ID != car.ID {
			t.Errorf("Expected ride to be associated with car ID %d, got %d", car.ID, fetchedRide.Car.ID)
		}

		if fetchedRide.Car.Name != "Tesla Model 3" {
			t.Errorf("Ride's associated car details are incorrect")
		}
	})

	t.Run("Fail to create ride without a car", func(t *testing.T) {
		description := "An invalid ride"
		ride := Rides{
			Title:       "Invalid Ride",
			Description: &description,
			BeginTime:   "2024-08-02T10:00:00Z",
			EndTime:     "2024-08-02T11:00:00Z",
			Category:    "Personal",
			// CarID is not set
		}

		result := db.Create(&ride)
		if result.Error == nil {
			t.Error("Expected an error when creating a ride without a car, but got none")
		}
	})
}

package db

import (
	"fahrtenbuch/pkg/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// initialize database connection
func Init(testing bool) {
	log.Println("Establishing Database connection...")

	username := os.Getenv("POSTGRES_USER")
	host := os.Getenv("POSTGRES_HOST")
	databaseName := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := "disable"

	if os.Getenv("IS_PROD") == "true" {
		sslmode = "require"
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		username, password, host, port, databaseName, sslmode)

	var err error
	if !testing {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}

	// check connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get raw database connection: %v", err)
	}

	// ping
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Auto migrate schema
	log.Println("Migarting Database Schemas")
	err = DB.AutoMigrate(
		&models.Organisation{},
		&models.User{},
		&models.Cars{},
		&models.Rides{},
	)
	if err != nil {
		log.Fatalf("Error migrating database: %s", err)
	}

	log.Println("Successfully connected to PostgreSQL")
}

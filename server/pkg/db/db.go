package db

import (
	"fahrtenbuch/pkg/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// initialize database connection
func Init() {
	username := os.Getenv("POSTGRES_USER")
	host := os.Getenv("POSTGRES_HOST")
	databaseName := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := "disable"

	if os.Getenv("IS_PROD") == "true" {
		sslmode = "require"
	}

	dns := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", username, password, host, port, databaseName, sslmode)

	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate schema
	err = DB.AutoMigrate(&models.User{}, &models.Organisation{}, &models.Rides{})
	if err != nil {
		log.Fatalf("Error migrating database: %s", err)
	}

	log.Println("Sucessfully connected to Postgres")
}

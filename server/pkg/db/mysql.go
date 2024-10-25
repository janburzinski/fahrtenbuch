package db

import (
	"log"
	"os"
	"server/pkg/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func Connect() error {
	//get mysql dsn from .env file
	dsn := os.Getenv("MYSQL_DSN")

	//connect to db
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//test db connection
	err = testConnection()
	if err != nil {
		return err
	}

	//auto migrate tables
	if err := db.AutoMigrate(&models.User{}, &models.Cars{}, &models.Rides{}); err != nil {
		panic(err)
	}

	return nil
}

func testConnection() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

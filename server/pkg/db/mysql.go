package db

import (
	"database/sql"
	"os"
	"server/pkg/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

var (
	db *bun.DB
)

func Connect() error {
	//connect to db
	dsn := os.Getenv("MYSQL_DSN")
	mysql, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db = bun.NewDB(mysql, mysqldialect.New())
	defer db.Close()

	//test db connection
	err = testConnection()
	if err != nil {
		return err
	}

	logger.Log("INFO", "Successfully connected to the MySQL Database!")

	return nil
}

func testConnection() error {
	err := db.Ping()
	return err
}

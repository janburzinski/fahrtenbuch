package db

import (
	"context"
	"database/sql"
	"os"
	"server/pkg/logger"
	"server/pkg/models"

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

	//run model migration
	err = runMigrations()
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

func runMigrations() error {
	ctx := context.Background()

	_, err := db.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	//todo: update and follow docs: https://bun.uptrace.dev/guide/starter-kit.html#app-structure
	return err
}

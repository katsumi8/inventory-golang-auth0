package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)
	if _, err := os.Stat("internal/storage/migrations"); os.IsNotExist(err) {
		log.Println("Migrations directory does not exist!")
	} else {
		log.Println("Migrations directory exists!")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/migrations", // マイグレーションファイルのパス
		"postgres",                           // DB名
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

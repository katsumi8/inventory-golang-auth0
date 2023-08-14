package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseConnection(db *sql.DB) error {
	return db.Close()
}

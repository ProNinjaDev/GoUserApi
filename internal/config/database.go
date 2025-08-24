package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDatabase() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "5432"
	}

	if user == "" {
		user = "myuser"
	}

	if password == "" {
		password = "12345"
	}

	if dbName == "" {
		dbName = "user_api_db"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("could not open the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	log.Println("Удалось подключиться к базе данных")
	return db, nil

}

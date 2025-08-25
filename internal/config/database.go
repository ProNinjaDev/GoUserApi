package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDatabase(cfg Config) (*sql.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName)

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

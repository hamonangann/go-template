package db

import (
	"database/sql"
	"fmt"
	"template/internal/config"
)

func ConnectDB() (*sql.DB, error) {
	DATABASE_USER := config.DATABASE_USER
	DATABASE_PASS := config.DATABASE_PASS
	DATABASE_HOST := config.DATABASE_HOST
	DATABASE_PORT := config.DATABASE_PORT
	DATABASE := config.DATABASE

	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DATABASE_USER, DATABASE_PASS, DATABASE_HOST, DATABASE_PORT, DATABASE)
	db, err := sql.Open("pgx", DATABASE_URL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

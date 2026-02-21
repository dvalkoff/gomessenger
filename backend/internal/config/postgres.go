package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB(dbConfig DbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.ConnectionStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

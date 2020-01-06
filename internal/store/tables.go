package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	ITEMS = `CREATE TABLE IF NOT EXISTS items (
		id bigserial not null primary key,
		title varchar not null,
		description varchar,
		date time not null,
		price float not null,
		images varchar[3]
	);`
)

func CreateTables(db *sql.DB) error {
	if _, err := db.Exec(ITEMS); err != nil {
		return err
	}
	return nil
}

func NewStore(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	if err := CreateTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

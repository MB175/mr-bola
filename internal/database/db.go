package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTable := `CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		owner TEXT NOT NULL,
		content TEXT NOT NULL
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Error creating notes table: %v", err)
	}

	return db, nil
}

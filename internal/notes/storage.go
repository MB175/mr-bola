package notes

import (
	"database/sql"
)

// InsertNote inserts a new note into the database
func InsertNote(db *sql.DB, note *Note) error {
	_, err := db.Exec("INSERT INTO notes (id, owner, content) VALUES (?, ?, ?)", note.ID, note.Owner, note.Content)
	return err
}

// GetNoteByID retrieves a note by its ID
func GetNoteByID(db *sql.DB, id string) (Note, error) {
	var note Note
	err := db.QueryRow("SELECT id, owner, content FROM notes WHERE id = ?", id).
		Scan(&note.ID, &note.Owner, &note.Content)
	return note, err
}

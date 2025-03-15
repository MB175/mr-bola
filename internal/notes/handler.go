package notes

import (
	"database/sql"
	"encoding/json"
	"github.com/MB175/mr-bola/internal/middleware"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

// CreateNoteHandler handles note creation
func CreateNoteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	claims, err := middleware.ExtractJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var note Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil || note.Content == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newUUID := uuid.New()
	note.ID = newUUID.String()
	note.Owner = claims.Subject // Set owner from JWT
	err = InsertNote(db, &note)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetNoteHandler retrieves a note by ID
func GetNoteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := strings.TrimPrefix(r.URL.Path, "/notes/")
	note, err := GetNoteByID(db, id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

package main

import (
	"github.com/MB175/mr-bola/internal/auth"
	"github.com/MB175/mr-bola/internal/database"
	"github.com/MB175/mr-bola/internal/notes"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := database.InitDB("notes.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/auth", auth.AuthHandler)
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			notes.CreateNoteHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			notes.GetNoteHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8071"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

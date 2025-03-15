package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MB175/mr-bola/internal/auth"
	"github.com/MB175/mr-bola/internal/database"
	"github.com/MB175/mr-bola/internal/notes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	// Initialize database
	db, err := database.InitDB("test_notes.db")
	if err != nil {
		panic("Failed to initialize test database")
	}

	// Create test server using an HTTP multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", auth.AuthHandler)
	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			notes.CreateNoteHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			notes.GetNoteHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start test HTTP server
	server = httptest.NewServer(mux)
	defer server.Close()

	// Run tests
	m.Run()
}

func getAuthToken(serverURL, username string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{"username": username})
	resp, err := http.Post(serverURL+"/auth", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["token"], nil
}

func createNote(serverURL, token, content string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{"content": content})
	req, err := http.NewRequest("POST", serverURL+"/notes", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status is not 201 Created
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create note: status %d, response: %s", resp.StatusCode, string(body))
	}

	// Decode JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Ensure "id" exists in the response
	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("error extracting note ID, response: %v", result)
	}

	return id, nil
}

func getNote(serverURL, token, noteID string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", serverURL+"/notes/"+noteID, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	return client.Do(req)
}

func TestObjectLevelAuthorization(t *testing.T) {
	serverURL := server.URL // Use dynamic server URL

	// User A gets a token
	tokenA, err := getAuthToken(serverURL, "userA")
	if err != nil {
		t.Fatalf("Failed to get token for userA: %v", err)
	}

	// User B gets a token
	tokenB, err := getAuthToken(serverURL, "userB")
	if err != nil {
		t.Fatalf("Failed to get token for userB: %v", err)
	}

	// User A creates a note
	noteID, err := createNote(serverURL, tokenA, "This is userA's secret note")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	// ✅ User A should be able to access their own note
	respA, err := getNote(serverURL, tokenA, noteID)
	if err != nil {
		t.Fatalf("Failed to fetch note for userA: %v", err)
	}
	if respA.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for userA, got %d", respA.StatusCode)
	}

	// ❌ User B should NOT be able to access user A's note
	respB, err := getNote(serverURL, tokenB, noteID)
	if err != nil {
		t.Fatalf("Failed to fetch note for userB: %v", err)
	}
	if respB.StatusCode == http.StatusOK {
		t.Fatalf("SECURITY BUG: User B was able to access User A's note!")
	}
	if respB.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status 403 for userB, got %d", respB.StatusCode)
	}

	// ❌ Anonymous user should NOT be able to access the note
	respAnon, err := getNote(serverURL, "", noteID)
	if err != nil {
		t.Fatalf("Failed to fetch note for anonymous user: %v", err)
	}
	if respAnon.StatusCode == http.StatusOK {
		t.Fatalf("SECURITY BUG: Anonymous user was able to access User A's note!")
	}
	if respAnon.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 for anonymous user, got %d", respAnon.StatusCode)
	}
}

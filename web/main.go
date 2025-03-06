package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sqlight/pkg/db"
	"sqlight/pkg/interfaces"
	"sqlight/pkg/sql"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message,omitempty"`
	Records []*interfaces.Record `json:"records,omitempty"`
}

const dbFile = "database.json"

func main() {
	// Create a new database
	database := db.NewDatabase()

	// Load existing database if it exists
	if _, err := os.Stat(dbFile); err == nil {
		if err := database.Load(dbFile); err != nil {
			log.Printf("Failed to load database: %v", err)
		}
	}

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/static/index.html")
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse and execute the query
		stmt, err := sql.ParseSQL(req.Query)
		if err != nil {
			json.NewEncoder(w).Encode(QueryResponse{
				Success: false,
				Message: fmt.Sprintf("Parse error: %v", err),
			})
			return
		}

		result, err := database.Execute(stmt)
		if err != nil {
			json.NewEncoder(w).Encode(QueryResponse{
				Success: false,
				Message: fmt.Sprintf("Execution error: %v", err),
			})
			return
		}

		// Save database after successful execution
		if err := database.Save(dbFile); err != nil {
			log.Printf("Failed to save database: %v", err)
		}

		json.NewEncoder(w).Encode(QueryResponse{
			Success: true,
			Message: result.Message,
			Records: result.Records,
		})
	})

	// Get list of tables
	http.HandleFunc("/tables", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tables := database.GetTables()
		json.NewEncoder(w).Encode(tables)
	})

	// Ensure the database file directory exists
	if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
		log.Fatal(err)
	}

	// Start server
	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

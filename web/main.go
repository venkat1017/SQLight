package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sqlight/pkg/db"
	"sqlight/pkg/sql"

	"github.com/gorilla/mux"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Success bool                   `json:"success"`
	Message string                `json:"message,omitempty"`
	Records []map[string]interface{} `json:"records,omitempty"`
	Columns []string              `json:"columns,omitempty"`
}

const dbFile = "database.json"

func main() {
	// Load database
	database, err := db.NewDatabase(dbFile)
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	// Create router
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Handle root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	})

	// Handle query execution
	r.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(QueryResponse{
				Success: false,
				Message: fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		// Parse and execute the query
		stmt, err := sql.Parse(req.Query)
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

		// Convert records to map format and get sorted columns
		var records []map[string]interface{}
		var columns []string

		// Process records for SELECT queries
		if result.Records != nil && len(result.Records) > 0 {
			// Use columns from the result
			columns = result.Columns
			if len(columns) == 0 && len(result.Records) > 0 {
				// If columns not provided, get them from the first record
				for col := range result.Records[0].Columns {
					columns = append(columns, col)
				}
				sort.Strings(columns)
			}

			// Convert records to maps
			for _, record := range result.Records {
				recordMap := make(map[string]interface{})
				for _, col := range columns {
					recordMap[col] = record.Columns[col]
				}
				records = append(records, recordMap)
			}

			// Log for debugging
			log.Printf("Query: %s", req.Query)
			log.Printf("Number of records: %d", len(records))
			log.Printf("Columns: %v", columns)
			if len(records) > 0 {
				log.Printf("First record: %+v", records[0])
			}
		}

		// Send response
		response := QueryResponse{
			Success:  true,
			Message:  result.Message,
			Records:  records,
			Columns:  columns,
		}

		// Log response for debugging
		log.Printf("Response: %+v", response)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	// Handle table list
	r.HandleFunc("/tables", func(w http.ResponseWriter, r *http.Request) {
		tables := database.GetTables()
		json.NewEncoder(w).Encode(map[string][]string{"tables": tables})
	})

	// Ensure the database file directory exists
	if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
		log.Fatal(err)
	}

	// Start server
	port := ":8081"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

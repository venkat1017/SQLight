package tests

import (
	"os"
	"sqlite-clone/pkg/db"
	"sqlite-clone/pkg/interfaces"
	"testing"
	"path/filepath"
)

func TestDatabase(t *testing.T) {
	// Create a temporary database file for testing
	tmpFile := "test_db.json"
	defer os.Remove(tmpFile)
	
	database := db.NewDatabase(tmpFile)

	// Test CREATE TABLE with data types
	err := database.CreateTable("users", []interfaces.ColumnDef{
		{Name: "id", Type: "INTEGER"},
		{Name: "name", Type: "TEXT"},
		{Name: "email", Type: "TEXT"},
	})
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Test duplicate table creation
	err = database.CreateTable("users", []interfaces.ColumnDef{
		{Name: "id", Type: "INTEGER"},
	})
	if err == nil {
		t.Error("Expected error when creating duplicate table")
	}

	// Test INSERT with data types
	r1 := interfaces.NewRecord(map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@email.com",
	})

	err = database.InsertIntoTable("users", r1)
	if err != nil {
		t.Fatalf("Error inserting record: %v", err)
	}

	// Test invalid table insert
	err = database.InsertIntoTable("nonexistent", r1)
	if err == nil {
		t.Error("Expected error when inserting into non-existent table")
	}

	// Test SELECT
	records, err := database.SelectFromTable("users")
	if err != nil {
		t.Fatalf("Error selecting from table: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	// Test invalid table select
	_, err = database.SelectFromTable("nonexistent")
	if err == nil {
		t.Error("Expected error when selecting from non-existent table")
	}

	// Test column values
	r, ok := records[0].(*interfaces.Record)
	if !ok {
		t.Fatalf("Invalid record type")
	}

	if r.Columns["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", r.Columns["name"])
	}

	// Test UPDATE
	err = database.UpdateTable("users", 
		map[string]interface{}{"name": "Alice Smith"}, 
		"id", 
		float64(1))
	if err != nil {
		t.Fatalf("Error updating record: %v", err)
	}

	// Verify update
	record, err := database.FindInTable("users", 1)
	if err != nil {
		t.Fatalf("Error finding updated record: %v", err)
	}
	r, _ = record.(*interfaces.Record)
	if r.Columns["name"] != "Alice Smith" {
		t.Errorf("Expected updated name 'Alice Smith', got '%v'", r.Columns["name"])
	}

	// Test invalid update
	err = database.UpdateTable("users", 
		map[string]interface{}{"name": "Bob"}, 
		"id", 
		float64(999))
	if err == nil {
		t.Error("Expected error when updating non-existent record")
	}

	// Test DELETE
	err = database.DeleteFromTable("users", "id", float64(1))
	if err != nil {
		t.Fatalf("Error deleting record: %v", err)
	}

	// Verify delete
	records, err = database.SelectFromTable("users")
	if err != nil {
		t.Fatalf("Error selecting after delete: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("Expected 0 records after delete, got %d", len(records))
	}

	// Test invalid delete
	err = database.DeleteFromTable("users", "id", float64(999))
	if err == nil {
		t.Error("Expected error when deleting non-existent record")
	}

	// Test file operations
	// Save current state
	err = database.Execute("INSERT INTO users VALUES (2, 'Bob', 'bob@email.com')")
	if err != nil {
		t.Fatalf("Error executing INSERT: %v", err)
	}

	// Create new database instance to load saved state
	database2 := db.NewDatabase(tmpFile)
	records, err = database2.SelectFromTable("users")
	if err != nil {
		t.Fatalf("Error selecting from loaded database: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("Expected 1 record in loaded database, got %d", len(records))
	}
}

func TestSQLParser(t *testing.T) {
	database := db.NewDatabase("test_parser.json")
	defer os.Remove("test_parser.json")

	// Test CREATE TABLE
	err := database.Execute("CREATE TABLE products (id INTEGER, name TEXT, price INTEGER)")
	if err != nil {
		t.Fatalf("Error executing CREATE TABLE: %v", err)
	}

	// Test INSERT
	err = database.Execute("INSERT INTO products VALUES (1, 'Widget', 100)")
	if err != nil {
		t.Fatalf("Error executing INSERT: %v", err)
	}

	// Test SELECT
	err = database.Execute("SELECT * FROM products")
	if err != nil {
		t.Fatalf("Error executing SELECT: %v", err)
	}

	// Test UPDATE
	err = database.Execute("UPDATE products SET price = 200 WHERE id = 1")
	if err != nil {
		t.Fatalf("Error executing UPDATE: %v", err)
	}

	// Test DELETE
	err = database.Execute("DELETE FROM products WHERE id = 1")
	if err != nil {
		t.Fatalf("Error executing DELETE: %v", err)
	}

	// Test invalid SQL
	err = database.Execute("INVALID SQL")
	if err == nil {
		t.Error("Expected error for invalid SQL")
	}

	err = database.Execute("SELECT * FROM nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

func TestBTree(t *testing.T) {
	tree := db.NewBTree()

	// Test Insert and Search
	r1 := interfaces.NewRecord(map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@email.com",
	})

	tree.Insert(1, r1)

	found := tree.Search(1)
	if found == nil {
		t.Fatal("Record not found after insertion")
	}

	if found.Columns["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", found.Columns["name"])
	}

	// Test non-existent key
	notFound := tree.Search(999)
	if notFound != nil {
		t.Error("Expected nil for non-existent key")
	}

	// Test multiple inserts
	r2 := interfaces.NewRecord(map[string]interface{}{
		"id":    2,
		"name":  "Bob",
		"email": "bob@email.com",
	})
	tree.Insert(2, r2)

	// Test Scan
	records := tree.Scan()
	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	// Test Delete
	tree.Delete(1)
	found = tree.Search(1)
	if found != nil {
		t.Error("Record still exists after deletion")
	}
}

func TestTransactions(t *testing.T) {
	// Create a temporary database file
	tmpFile := filepath.Join(os.TempDir(), "test_db_tx.json")
	defer os.Remove(tmpFile)

	db := db.NewDatabase(tmpFile)

	// Create a test table
	columns := []interfaces.ColumnDef{
		{Name: "id", Type: "INTEGER"},
		{Name: "name", Type: "TEXT"},
	}
	err := db.CreateTable("users", columns)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Test 1: Basic transaction commit
	t.Run("Transaction Commit", func(t *testing.T) {
		err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert a record
		record := interfaces.NewRecord(map[string]interface{}{
			"id":   1,
			"name": "Alice",
		})
		tables := db.Tables()
		err = tables["users"].Insert(record)
		if err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		err = db.Commit()
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify record exists
		tables = db.Tables()
		cursor := tables["users"].NewCursor()
		found := false
		record, err = cursor.First()
		for record != nil && err == nil {
			if record.Columns["id"] == 1 {
				found = true
				break
			}
			record, err = cursor.Next()
		}
		if !found {
			t.Error("Record not found after commit")
		}
	})

	// Test 2: Transaction rollback
	t.Run("Transaction Rollback", func(t *testing.T) {
		err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert a record
		record := interfaces.NewRecord(map[string]interface{}{
			"id":   2,
			"name": "Bob",
		})
		tables := db.Tables()
		err = tables["users"].Insert(record)
		if err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		err = db.Rollback()
		if err != nil {
			t.Fatalf("Failed to rollback transaction: %v", err)
		}

		// Verify record does not exist
		tables = db.Tables()
		cursor := tables["users"].NewCursor()
		found := false
		record, err = cursor.First()
		for record != nil && err == nil {
			if record.Columns["id"] == 2 {
				found = true
				break
			}
			record, err = cursor.Next()
		}
		if found {
			t.Error("Record found after rollback")
		}
	})

	// Test 3: Nested transactions not allowed
	t.Run("Nested Transactions", func(t *testing.T) {
		err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin first transaction: %v", err)
		}

		err = db.Begin()
		if err == nil {
			t.Error("Expected error when beginning nested transaction")
		}

		err = db.Rollback()
		if err != nil {
			t.Fatalf("Failed to rollback transaction: %v", err)
		}
	})
}

package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sqlight/pkg/interfaces"
	"sqlight/pkg/sql"
)

// Verify Database implements the interface
var _ interfaces.Database = (*Database)(nil)

// Statement defines the interface for SQL statements
type Statement interface {
	Exec(db *Database) error
}

// SelectStatement represents a SELECT query
type SelectStatement struct {
	Table string
}

func (s *SelectStatement) Exec(db *Database) error {
	return nil // Placeholder for Step 3
}

// InsertStatement represents an INSERT command
type InsertStatement struct {
	Table  string
	Values []interface{}
}

func (s *InsertStatement) Exec(db *Database) error {
	return nil // Placeholder for Step 3
}

// CreateStatement represents a CREATE TABLE command
type CreateStatement struct {
	Table   string
	Columns []string
}

func (s *CreateStatement) Exec(db *Database) error {
	return nil // Placeholder for Step 3
}

// Database holds the tables
type Database struct {
	tables        map[string]*Table
	file          string // Path to the database file
	inTransaction bool
	currentTx     *Transaction
}

// Tables returns a copy of the tables map
func (db *Database) Tables() map[string]*Table {
	fmt.Printf("Getting tables: %v\n", db.tables)
	return db.tables
}

// SetTables allows setting the tables map directly
func (db *Database) SetTables(tables map[string]*Table) {
	db.tables = tables
}

// NewDatabase initializes a new Database and loads existing data from the file
func NewDatabase(file string) *Database {
	db := &Database{
		tables: make(map[string]*Table),
		file:   file,
	}
	db.load() // Load existing data from the file
	return db
}

// Load data from the file
func (db *Database) load() {
	if _, err := os.Stat(db.file); err == nil {
		file, err := os.Open(db.file)
		if err != nil {
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&db.tables)
		if err != nil {
			return
		}
	}
}

// Save data to the file
func (db *Database) Save() {
	fmt.Println("Saving database state:", db.tables) // Debugging output

	file, err := os.Create(db.file)
	if err != nil {
		fmt.Println("Error creating database file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(db.tables)
	if err != nil {
		fmt.Println("Error encoding database to JSON:", err)
		return
	}
}

// Begin starts a new transaction
func (db *Database) Begin() error {
	if db.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}

	// Create a new transaction
	tx := NewTransaction()
	db.currentTx = tx

	// Create a snapshot of the current state
	db.currentTx.CreateSnapshot(db.tables)
	db.inTransaction = true
	return nil
}

// Commit commits the current transaction
func (db *Database) Commit() error {
	if !db.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if db.currentTx == nil {
		return fmt.Errorf("no active transaction")
	}

	// Save the changes to disk
	db.Save()
	db.currentTx = nil
	db.inTransaction = false
	return nil
}

// Rollback rolls back the current transaction
func (db *Database) Rollback() error {
	if !db.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if db.currentTx == nil {
		return fmt.Errorf("no active transaction")
	}

	// Restore the database state from the snapshot
	for name, table := range db.currentTx.snapshot {
		db.tables[name] = table.Clone()
	}

	// Handle deleted tables
	for name := range db.currentTx.deleted {
		delete(db.tables, name)
	}

	db.currentTx = nil
	db.inTransaction = false
	return nil
}

// CreateTable creates a new table with the specified columns
func (db *Database) CreateTable(name string, columns []interfaces.ColumnDef) error {
	fmt.Printf("Creating table: %s with columns: %v\n", name, columns)

	// Check if table already exists
	if _, exists := db.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}

	table, err := NewTable(columns)
	if err != nil {
		return err
	}
	db.tables[name] = table
	return nil
}

func (db *Database) GetTableColumns(name string) ([]string, error) {
	if table, ok := db.tables[name]; ok {
		return table.Columns(), nil
	}
	return nil, fmt.Errorf("table %s not found", name)
}

// InsertIntoTable inserts a record into a table
func (db *Database) InsertIntoTable(name string, record *interfaces.Record) error {
	table, exists := db.tables[name]
	if !exists {
		return fmt.Errorf("table %s does not exist", name)
	}

	return table.Insert(record)
}

// SelectFromTable selects records from a table
func (db *Database) SelectFromTable(name string, whereColumn string, whereValue interface{}) ([]interface{}, error) {
	table, exists := db.tables[name]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", name)
	}

	records := table.Select()
	result := make([]interface{}, len(records))
	for i, r := range records {
		result[i] = r
	}

	if whereColumn == "" {
		return result, nil
	}

	// Filter records based on where clause
	var filtered []interface{}
	for _, record := range result {
		if r, ok := record.(*interfaces.Record); ok {
			if val, exists := r.Columns[whereColumn]; exists {
				// Convert both values to strings for comparison
				valStr := fmt.Sprintf("%v", val)
				whereStr := fmt.Sprintf("%v", whereValue)
				if valStr == whereStr {
					filtered = append(filtered, record)
				}
			}
		}
	}

	return filtered, nil
}

// Add new method for finding specific records
func (db *Database) FindInTable(name string, id int) (interface{}, error) {
	if table, ok := db.tables[name]; ok {
		record := table.Find(id)
		if record == nil {
			return nil, fmt.Errorf("record with id %d not found", id)
		}
		return record, nil
	}
	return nil, fmt.Errorf("table %s not found", name)
}

// DeleteFromTable deletes records from a table that match the where clause
func (db *Database) DeleteFromTable(name string, whereCol string, whereVal interface{}) error {
	if table, ok := db.tables[name]; ok {
		return table.Delete(whereCol, whereVal)
	}
	return fmt.Errorf("table %s not found", name)
}

// UpdateTable updates records in a table that match the where clause
func (db *Database) UpdateTable(name string, setColumns map[string]interface{}, whereCol string, whereVal interface{}) error {
	if table, ok := db.tables[name]; ok {
		return table.Update(setColumns, whereCol, whereVal)
	}
	return fmt.Errorf("table %s not found", name)
}

// Update Execute to match the interface
func (db *Database) Execute(query string) error {
	stmt, err := sql.ParseSQL(query)
	if err != nil {
		return err
	}
	if stmt != nil {
		err = stmt.Exec(db)
		if err != nil {
			return err
		}
		// Save changes after successful execution
		db.Save()
	}
	return nil
}

func (db *Database) PrintState() {
	for tableName, table := range db.tables {
		fmt.Printf("Table: %s\n", tableName)
		records := table.Select()
		for _, record := range records {
			fmt.Printf("Record: %+v\n", record)
		}
	}
}

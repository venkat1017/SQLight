package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sqlight/pkg/interfaces"
	"sync"
)

// Database represents a SQLite database
type Database struct {
	tables map[string]*Table
	mutex  sync.RWMutex
}

// DatabaseSnapshot represents the database state for persistence
type DatabaseSnapshot struct {
	Tables map[string]TableSnapshot `json:"tables"`
}

// TableSnapshot represents a table's state for persistence
type TableSnapshot struct {
	Name    string                         `json:"name"`
	Columns []interfaces.ColumnDef         `json:"columns"`
	Records map[string]map[string]interface{} `json:"records"`
}

// NewDatabase creates a new database instance
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

// Execute executes a SQL statement
func (d *Database) Execute(stmt interfaces.Statement) (*interfaces.Result, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	switch s := stmt.(type) {
	case *interfaces.CreateStatement:
		return d.executeCreate(s)
	case *interfaces.InsertStatement:
		return d.executeInsert(s)
	case *interfaces.SelectStatement:
		return d.executeSelect(s)
	default:
		return nil, fmt.Errorf("unsupported statement type")
	}
}

// executeCreate handles CREATE TABLE statements
func (d *Database) executeCreate(stmt *interfaces.CreateStatement) (*interfaces.Result, error) {
	// Check if table already exists
	if _, exists := d.tables[stmt.TableName]; exists {
		return nil, fmt.Errorf("table %s already exists", stmt.TableName)
	}

	// Create new table
	table, err := NewTable(stmt.Columns)
	if err != nil {
		return nil, err
	}

	// Set table name
	table.SetName(stmt.TableName)

	// Add table to database
	d.tables[stmt.TableName] = table

	return &interfaces.Result{
		Success: true,
		Message: fmt.Sprintf("Table %s created successfully", stmt.TableName),
	}, nil
}

// executeInsert handles INSERT statements
func (d *Database) executeInsert(stmt *interfaces.InsertStatement) (*interfaces.Result, error) {
	// Get table
	table, exists := d.tables[stmt.TableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	// Create record from values
	record := &interfaces.Record{
		Columns: stmt.Values,
	}

	// Insert record
	err := table.Insert(record)
	if err != nil {
		return nil, err
	}

	return &interfaces.Result{
		Success: true,
		Message: "Record inserted successfully",
	}, nil
}

// executeSelect handles SELECT statements
func (d *Database) executeSelect(stmt *interfaces.SelectStatement) (*interfaces.Result, error) {
	// Get table
	table, exists := d.tables[stmt.TableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	// Get all records for now (we'll add filtering later)
	records := table.GetRecords()

	return &interfaces.Result{
		Success: true,
		Records: records,
	}, nil
}

// GetTables returns all table names in the database
func (d *Database) GetTables() []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	tables := make([]string, 0, len(d.tables))
	for name := range d.tables {
		tables = append(tables, name)
	}
	return tables
}

// GetTable returns a specific table by name
func (d *Database) GetTable(name string) (interfaces.Table, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	table, exists := d.tables[name]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", name)
	}
	return table, nil
}

// Save saves the database to a file
func (d *Database) Save(filename string) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Create database snapshot
	snapshot := DatabaseSnapshot{
		Tables: make(map[string]TableSnapshot),
	}

	for name, table := range d.tables {
		records := table.GetRecords()
		recordMap := make(map[string]map[string]interface{})
		
		for i, record := range records {
			recordMap[fmt.Sprintf("record_%d", i)] = record.Columns
		}

		snapshot.Tables[name] = TableSnapshot{
			Name:    name,
			Columns: table.GetColumnDefs(),
			Records: recordMap,
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write database file: %v", err)
	}

	return nil
}

// Load loads the database from a file
func (d *Database) Load(filename string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read database file: %v", err)
	}

	// Unmarshal JSON
	var snapshot DatabaseSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal database: %v", err)
	}

	// Restore tables
	d.tables = make(map[string]*Table)
	for name, tableSnapshot := range snapshot.Tables {
		// Create table
		table, err := NewTable(tableSnapshot.Columns)
		if err != nil {
			return fmt.Errorf("failed to create table %s: %v", name, err)
		}

		// Set table name
		table.SetName(name)

		// Insert records
		for _, record := range tableSnapshot.Records {
			if err := table.Insert(&interfaces.Record{Columns: record}); err != nil {
				return fmt.Errorf("failed to restore record in table %s: %v", name, err)
			}
		}

		d.tables[name] = table
	}

	return nil
}

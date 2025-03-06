package db

import (
	"fmt"
	"sqlight/pkg/interfaces"
	"sync"
)

// Table represents a database table
type Table struct {
	name     string
	columns  []interfaces.ColumnDef
	records  []*interfaces.Record
	mutex    sync.RWMutex
}

// NewTable creates a new table with the given column definitions
func NewTable(columns []interfaces.ColumnDef) (*Table, error) {
	if len(columns) == 0 {
		return nil, fmt.Errorf("table must have at least one column")
	}

	return &Table{
		columns: columns,
		records: make([]*interfaces.Record, 0),
	}, nil
}

// GetName returns the table name
func (t *Table) GetName() string {
	return t.name
}

// SetName sets the table name
func (t *Table) SetName(name string) {
	t.name = name
}

// Clone creates a deep copy of the table
func (t *Table) Clone() interfaces.Table {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	newTable := &Table{
		name:    t.name,
		columns: make([]interfaces.ColumnDef, len(t.columns)),
		records: make([]*interfaces.Record, len(t.records)),
	}

	// Copy columns
	copy(newTable.columns, t.columns)

	// Copy records
	for i, record := range t.records {
		newRecord := &interfaces.Record{
			Columns: make(map[string]interface{}),
		}
		for k, v := range record.Columns {
			newRecord.Columns[k] = v
		}
		newTable.records[i] = newRecord
	}

	return newTable
}

// Insert adds a new record to the table
func (t *Table) Insert(record *interfaces.Record) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Validate record against column definitions
	if err := t.validateRecord(record); err != nil {
		return err
	}

	// Add record
	t.records = append(t.records, record)
	return nil
}

// GetRecords returns all records in the table
func (t *Table) GetRecords() []*interfaces.Record {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Return a copy of records to prevent modification
	records := make([]*interfaces.Record, len(t.records))
	copy(records, t.records)
	return records
}

// GetColumns returns the column names for the table
func (t *Table) GetColumns() []string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	columns := make([]string, len(t.columns))
	for i, col := range t.columns {
		columns[i] = col.Name
	}
	return columns
}

// GetColumnDefs returns the column definitions for the table
func (t *Table) GetColumnDefs() []interfaces.ColumnDef {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Return a copy of column definitions
	columns := make([]interfaces.ColumnDef, len(t.columns))
	copy(columns, t.columns)
	return columns
}

// validateRecord checks if a record matches the table's schema
func (t *Table) validateRecord(record *interfaces.Record) error {
	if record == nil || record.Columns == nil {
		return fmt.Errorf("invalid record: record or columns are nil")
	}

	// Check that all required columns are present
	for _, col := range t.columns {
		if _, exists := record.Columns[col.Name]; !exists {
			return fmt.Errorf("missing required column: %s", col.Name)
		}
	}

	return nil
}

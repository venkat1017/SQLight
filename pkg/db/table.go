package db

import (
	"fmt"
	"sqlight/pkg/interfaces"
)

// Table represents a database table
type Table struct {
	name    string
	columns []interfaces.Column
	records []*interfaces.Record
}

// NewTable creates a new table with the given columns
func NewTable(columns []interfaces.Column) (*Table, error) {
	if len(columns) == 0 {
		return nil, fmt.Errorf("table must have at least one column")
	}

	return &Table{
		columns: columns,
		records: make([]*interfaces.Record, 0),
	}, nil
}

// SetName sets the table name
func (t *Table) SetName(name string) {
	t.name = name
}

// GetName returns the table name
func (t *Table) GetName() string {
	return t.name
}

// Insert adds a new record to the table
func (t *Table) Insert(record *interfaces.Record) error {
	// Validate record against column definitions
	for _, col := range t.columns {
		value, exists := record.Columns[col.Name]

		// Check NOT NULL constraint
		if !col.Nullable && (!exists || value == nil) {
			return fmt.Errorf("column %s cannot be null", col.Name)
		}

		// Check PRIMARY KEY and UNIQUE constraints
		if (col.PrimaryKey || col.Unique) && exists {
			for _, existingRecord := range t.records {
				if existingValue, ok := existingRecord.Columns[col.Name]; ok && existingValue == value {
					constraint := "UNIQUE"
					if col.PrimaryKey {
						constraint = "PRIMARY KEY"
					}
					return fmt.Errorf("duplicate value in %s column %s", constraint, col.Name)
				}
			}
		}
	}

	t.records = append(t.records, record)
	return nil
}

// GetRecords returns all records in the table
func (t *Table) GetRecords() []*interfaces.Record {
	return t.records
}

// GetColumns returns all column names
func (t *Table) GetColumns() []string {
	names := make([]string, len(t.columns))
	for i, col := range t.columns {
		names[i] = col.Name
	}
	return names
}

// GetColumnDefs returns all column definitions
func (t *Table) GetColumnDefs() []interfaces.Column {
	return t.columns
}

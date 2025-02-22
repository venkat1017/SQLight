package interfaces

// ColumnDef represents a column definition in a table
type ColumnDef struct {
	Name string
	Type string // The data type as a string (e.g., "INTEGER", "TEXT")
}

// Record represents a database record with dynamic columns
type Record struct {
	Columns map[string]interface{} `json:"columns"`
}

// NewRecord creates a new record with the given columns
func NewRecord(values map[string]interface{}) *Record {
	return &Record{
		Columns: values,
	}
}

// Statement interface for SQL commands
type Statement interface {
	Exec(db Database) error
}

// Database interface for SQL operations
type Database interface {
	// Table operations
	CreateTable(name string, columns []ColumnDef) error
	GetTableColumns(name string) ([]string, error)
	InsertIntoTable(name string, record *Record) error
	UpdateTable(name string, setColumns map[string]interface{}, whereColumn string, whereValue interface{}) error
	DeleteFromTable(name string, whereColumn string, whereValue interface{}) error
	SelectFromTable(name string, whereColumn string, whereValue interface{}) ([]interface{}, error)

	// Transaction operations
	Begin() error
	Commit() error
	Rollback() error
}

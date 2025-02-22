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

// Database interface for SQL operations
type Database interface {
	CreateTable(name string, columns []ColumnDef) error
	InsertIntoTable(name string, record interface{}) error
	SelectFromTable(name string) ([]interface{}, error)
	FindInTable(name string, id int) (interface{}, error)
	GetTableColumns(name string) ([]string, error)
	DeleteFromTable(name string, whereCol string, whereVal interface{}) error
	UpdateTable(name string, setColumns map[string]interface{}, whereCol string, whereVal interface{}) error
}

// Statement interface for SQL commands
type Statement interface {
	Exec(db Database) error
}

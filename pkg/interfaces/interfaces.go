package interfaces

// Statement represents a SQL statement
type Statement interface {
	Type() string
}

// CreateStatement represents a CREATE TABLE statement
type CreateStatement struct {
	TableName string
	Columns   []ColumnDef
}

func (s *CreateStatement) Type() string {
	return "CREATE"
}

// InsertStatement represents an INSERT statement
type InsertStatement struct {
	TableName string
	Values    map[string]interface{}
}

func (s *InsertStatement) Type() string {
	return "INSERT"
}

// SelectStatement represents a SELECT statement
type SelectStatement struct {
	TableName string
	Columns   []string
	Where     string
}

func (s *SelectStatement) Type() string {
	return "SELECT"
}

// Record represents a database record
type Record struct {
	Columns map[string]interface{}
}

// Result represents the result of executing a statement
type Result struct {
	Success bool
	Message string
	Records []*Record
}

// ColumnDef represents a column definition
type ColumnDef struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	PrimaryKey bool   `json:"primary_key,omitempty"`
	NotNull    bool   `json:"not_null,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	References string `json:"references,omitempty"`
}

// Table represents a database table
type Table interface {
	GetName() string
	Insert(record *Record) error
	GetRecords() []*Record
	GetColumns() []string
	GetColumnDefs() []ColumnDef
	Clone() Table
}

// Database represents a database
type Database interface {
	Execute(stmt Statement) (*Result, error)
	GetTables() []string
	GetTable(name string) (Table, error)
	Save(filename string) error
	Load(filename string) error
}

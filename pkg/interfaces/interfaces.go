package interfaces

// Statement represents a SQL statement
type Statement interface {
	Type() string
}

// Column represents a table column definition
type Column struct {
	Name       string
	Type       string
	PrimaryKey bool
	Nullable   bool
	Unique     bool
}

// Table represents a database table
type Table struct {
	Name    string
	Columns []Column
	Records []*Record
}

// CreateStatement represents a CREATE TABLE statement
type CreateStatement struct {
	TableName string
	Columns   []Column
}

func (s *CreateStatement) Type() string {
	return "CREATE"
}

// InsertStatement represents an INSERT statement
type InsertStatement struct {
	TableName string
	Columns   []string
	Values    []interface{}
}

func (s *InsertStatement) Type() string {
	return "INSERT"
}

// SelectStatement represents a SELECT statement
type SelectStatement struct {
	TableName string
	Columns   []string
	Where     map[string]interface{}
}

func (s *SelectStatement) Type() string {
	return "SELECT"
}

// DropStatement represents a DROP TABLE statement
type DropStatement struct {
	TableName string
}

func (s *DropStatement) Type() string {
	return "DROP"
}

// DescribeStatement represents a DESCRIBE TABLE statement
type DescribeStatement struct {
	TableName string
}

func (s *DescribeStatement) Type() string {
	return "DESCRIBE"
}

// DeleteStatement represents a DELETE statement
type DeleteStatement struct {
	TableName string
	Where     map[string]interface{}
}

func (s *DeleteStatement) Type() string {
	return "DELETE"
}

// BeginTransactionStatement represents a BEGIN TRANSACTION statement
type BeginTransactionStatement struct{}

func (s *BeginTransactionStatement) Type() string {
	return "BEGIN TRANSACTION"
}

// CommitStatement represents a COMMIT statement
type CommitStatement struct{}

func (s *CommitStatement) Type() string {
	return "COMMIT"
}

// RollbackStatement represents a ROLLBACK statement
type RollbackStatement struct{}

func (s *RollbackStatement) Type() string {
	return "ROLLBACK"
}

// Record represents a database record
type Record struct {
	Columns map[string]interface{}
}

// Result represents a database operation result
type Result struct {
	Success  bool
	Message  string
	Records  []*Record
	Columns  []string
	IsSelect bool
}

// Transaction represents a database transaction
type Transaction struct {
	Tables map[string]*Table
}

// Database represents a database
type Database interface {
	Execute(stmt Statement) (*Result, error)
	GetTables() []string
	GetTable(name string) (Table, error)
	Save(filename string) error
	Load(filename string) error
	BeginTransaction() error
	Commit() error
	Rollback() error
}

package db

// Transaction represents a database transaction
type Transaction struct {
	snapshot map[string]*Table // Snapshot of tables at transaction start
	deleted  map[string]bool   // Tables marked for deletion
}

// NewTransaction creates a new transaction
func NewTransaction() *Transaction {
	return &Transaction{
		snapshot: make(map[string]*Table),
		deleted:  make(map[string]bool),
	}
}

// CreateSnapshot creates a snapshot of the current database state
func (tx *Transaction) CreateSnapshot(tables map[string]*Table) {
	tx.snapshot = make(map[string]*Table)
	for name, table := range tables {
		tx.snapshot[name] = table.Clone()
	}
}

// MarkForDeletion marks a table for deletion
func (tx *Transaction) MarkForDeletion(tableName string) {
	tx.deleted[tableName] = true
}

// IsMarkedForDeletion checks if a table is marked for deletion
func (tx *Transaction) IsMarkedForDeletion(tableName string) bool {
	return tx.deleted[tableName]
}

// Commit applies the changes in the transaction
func (tx *Transaction) Commit() error {
	return nil
}

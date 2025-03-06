package db

import (
	"sqlight/pkg/interfaces"
	"sync"
)

// Transaction represents a database transaction
type Transaction struct {
	db      *Database
	tables  map[string]*Table
	mutex   sync.RWMutex
	started bool
}

// NewTransaction creates a new transaction
func NewTransaction(db *Database) *Transaction {
	return &Transaction{
		db:     db,
		tables: make(map[string]*Table),
	}
}

// Begin starts the transaction
func (t *Transaction) Begin() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.started {
		return nil
	}

	t.started = true
	return nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.started {
		return nil
	}

	t.started = false
	return nil
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.started {
		return nil
	}

	t.started = false
	return nil
}

// Execute executes a statement within the transaction
func (t *Transaction) Execute(stmt interfaces.Statement) (*interfaces.Result, error) {
	return t.db.Execute(stmt)
}

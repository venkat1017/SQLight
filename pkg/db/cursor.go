package db

import "sqlight/pkg/interfaces"

// Cursor represents a database cursor
type Cursor struct {
	table  *Table
	pos    int
}

// NewCursor creates a new cursor for the table
func NewCursor(t *Table) *Cursor {
	return &Cursor{
		table: t,
		pos:   -1,
	}
}

// Next moves the cursor to the next record
func (c *Cursor) Next() bool {
	records := c.table.GetRecords()
	c.pos++
	return c.pos < len(records)
}

// Current returns the current record
func (c *Cursor) Current() *interfaces.Record {
	records := c.table.GetRecords()
	if c.pos >= 0 && c.pos < len(records) {
		return records[c.pos]
	}
	return nil
}

package db

import "sqlight/pkg/interfaces"

// Cursor represents a cursor for iterating over records in a table
type Cursor struct {
	table    *Table
	current  *interfaces.Record
	position int
	records  []*interfaces.Record
}

// NewCursor creates a new cursor for the table
func (t *Table) NewCursor() *Cursor {
	records := t.tree.Scan()
	cursor := &Cursor{
		table:    t,
		records:  records,
		position: -1,
	}
	return cursor
}

// First moves the cursor to the first record
func (c *Cursor) First() (*interfaces.Record, error) {
	if len(c.records) == 0 {
		return nil, nil
	}
	c.position = 0
	c.current = c.records[c.position]
	return c.current, nil
}

// Next moves the cursor to the next record
func (c *Cursor) Next() (*interfaces.Record, error) {
	if c.position >= len(c.records)-1 {
		return nil, nil
	}
	c.position++
	c.current = c.records[c.position]
	return c.current, nil
}

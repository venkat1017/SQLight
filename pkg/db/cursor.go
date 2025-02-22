package db

import "sqlite-clone/pkg/interfaces"

type Cursor struct {
	table   *Table
	current int
	records []*interfaces.Record // Cache records for iteration
}

func (t *Table) NewCursor() *Cursor {
	return &Cursor{
		table:   t,
		current: -1,
		records: t.Select(), // Get all records from B-Tree
	}
}

func (c *Cursor) Next() *interfaces.Record {
	c.current++
	if c.current < len(c.records) {
		return c.records[c.current]
	}
	return nil
}

func (c *Cursor) Reset() {
	c.current = -1
}

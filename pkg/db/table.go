package db

import (
	"encoding/json"
	"fmt"
	"sqlite-clone/pkg/interfaces"
	"sqlite-clone/pkg/types/datatypes"
	"sqlite-clone/pkg/logger"
)

type Column struct {
	Name string
	Type datatypes.DataType
}

// MarshalJSON customizes the JSON serialization of Column
func (c Column) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Name: c.Name,
		Type: c.Type.Name(),
	})
}

// UnmarshalJSON customizes the JSON deserialization of Column
func (c *Column) UnmarshalJSON(data []byte) error {
	var temp struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	c.Name = temp.Name
	dataType, err := datatypes.GetType(temp.Type)
	if err != nil {
		return err
	}
	c.Type = dataType
	return nil
}

type Table struct {
	columns []Column
	tree    *BTree
}

func NewTable(columnDefs []interfaces.ColumnDef) (*Table, error) {
	columns := make([]Column, len(columnDefs))
	for i, def := range columnDefs {
		dataType, err := datatypes.GetType(def.Type)
		if err != nil {
			return nil, fmt.Errorf("column %s: %v", def.Name, err)
		}
		columns[i] = Column{
			Name: def.Name,
			Type: dataType,
		}
	}

	return &Table{
		columns: columns,
		tree:    NewBTree(),
	}, nil
}

func (t *Table) Columns() []string {
	columnNames := make([]string, len(t.columns))
	for i, col := range t.columns {
		columnNames[i] = col.Name
	}
	return columnNames
}

func (t *Table) Insert(r *interfaces.Record) error {
	// Validate all columns
	for _, col := range t.columns {
		value, ok := r.Columns[col.Name]
		if !ok {
			return fmt.Errorf("missing column: %s", col.Name)
		}

		// Validate and convert the value
		if err := col.Type.Validate(value); err != nil {
			return fmt.Errorf("column %s: %v", col.Name, err)
		}

		// Convert the value to the correct type
		converted, err := col.Type.Convert(value)
		if err != nil {
			return fmt.Errorf("column %s: %v", col.Name, err)
		}

		// Update the record with the converted value
		r.Columns[col.Name] = converted
	}

	// Get the ID
	idVal, ok := r.Columns["id"]
	if !ok {
		return fmt.Errorf("record must have an 'id' column")
	}

	// Handle both float64 and int for ID
	var id int
	switch v := idVal.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return fmt.Errorf("'id' must be a number")
	}

	// Insert the record into the B-tree
	err := t.tree.Insert(id, r)
	if err != nil {
		return fmt.Errorf("failed to insert record into B-tree: %v", err)
	}

	return nil
}

func (t *Table) Select() []*interfaces.Record {
	return t.tree.Scan()
}

func (t *Table) Find(id int) *interfaces.Record {
	return t.tree.Search(id)
}

func (t *Table) Delete(whereCol string, whereVal interface{}) error {
	logger.Debugf("Delete called with whereCol=%s, whereVal=%v (type: %T)", whereCol, whereVal, whereVal)

	// If no where clause, delete all records
	if whereCol == "" {
		t.tree = NewBTree()
		return nil
	}

	// Get all records
	records := t.tree.Scan()
	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	// Create a new B-tree for the remaining records
	newTree := NewBTree()
	deleted := false

	// Copy records that don't match the where clause
	for _, record := range records {
		if val, ok := record.Columns[whereCol]; ok {
			logger.Debugf("Comparing record value=%v (type: %T) with whereVal=%v (type: %T)", val, val, whereVal, whereVal)
			isMatch := false

			// Convert both values to float64 for comparison if they're numbers
			if whereNum, ok := whereVal.(float64); ok {
				logger.Debugf("whereVal is float64: %v", whereNum)
				switch v := val.(type) {
				case float64:
					isMatch = v == whereNum
					logger.Debugf("val is float64: %v, isMatch=%v", v, isMatch)
				case int:
					isMatch = float64(v) == whereNum
					logger.Debugf("val is int: %v, isMatch=%v", v, isMatch)
				default:
					logger.Debugf("val is neither float64 nor int: %v", v)
				}
			} else if whereStr, ok := whereVal.(string); ok {
				logger.Debugf("whereVal is string: %v", whereStr)
				if valStr, ok := val.(string); ok {
					isMatch = valStr == whereStr
					logger.Debugf("val is string: %v, isMatch=%v", valStr, isMatch)
				}
			}

			if !isMatch {
				// Keep this record (not matching delete criteria)
				var id int
				switch v := record.Columns["id"].(type) {
				case float64:
					id = int(v)
				case int:
					id = v
				default:
					return fmt.Errorf("invalid id type: %T", v)
				}
				newTree.Insert(id, record)
			} else {
				logger.Debugf("Match found! Deleting record")
				deleted = true
			}
		}
	}

	if !deleted {
		return fmt.Errorf("no records matched the delete criteria")
	}

	t.tree = newTree
	return nil
}

func (t *Table) Update(setColumns map[string]interface{}, whereCol string, whereVal interface{}) error {
	logger.Debugf("Update called with whereCol=%s, whereVal=%v (type: %T)", whereCol, whereVal, whereVal)

	// Get all records
	records := t.tree.Scan()
	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	// Create a new B-tree for the updated records
	newTree := NewBTree()
	updated := false

	// Update matching records
	for _, record := range records {
		if val, ok := record.Columns[whereCol]; ok {
			logger.Debugf("Comparing record value=%v (type: %T) with whereVal=%v (type: %T)", val, val, whereVal, whereVal)
			isMatch := false
			
			// Convert both values to float64 for comparison if they're numbers
			if whereNum, ok := whereVal.(float64); ok {
				logger.Debugf("whereVal is float64: %v", whereNum)
				switch v := val.(type) {
				case float64:
					isMatch = v == whereNum
					logger.Debugf("val is float64: %v, isMatch=%v", v, isMatch)
				case int:
					isMatch = float64(v) == whereNum
					logger.Debugf("val is int: %v, isMatch=%v", v, isMatch)
				default:
					logger.Debugf("val is neither float64 nor int: %v", v)
				}
			} else if whereStr, ok := whereVal.(string); ok {
				logger.Debugf("whereVal is string: %v", whereStr)
				if valStr, ok := val.(string); ok {
					isMatch = valStr == whereStr
					logger.Debugf("val is string: %v, isMatch=%v", valStr, isMatch)
				}
			}

			if isMatch {
				logger.Debugf("Match found! Updating record")
				// Update the matching record
				for k, v := range setColumns {
					record.Columns[k] = v
				}
				updated = true
			}

			// Keep the record (either updated or unchanged)
			var id int
			switch v := record.Columns["id"].(type) {
			case float64:
				id = int(v)
			case int:
				id = v
			default:
				return fmt.Errorf("invalid id type: %T", v)
			}
			newTree.Insert(id, record)
		}
	}

	if !updated {
		return fmt.Errorf("no records matched the update criteria")
	}

	t.tree = newTree
	return nil
}

// Clone creates a deep copy of the table
func (t *Table) Clone() *Table {
	newTable := &Table{
		columns: make([]Column, len(t.columns)),
		tree:    NewBTree(),
	}

	// Copy columns
	copy(newTable.columns, t.columns)

	// Copy records
	records := t.tree.Scan()
	for _, record := range records {
		// Create a new map for the record's columns
		newColumns := make(map[string]interface{})
		for k, v := range record.Columns {
			newColumns[k] = v
		}
		newRecord := interfaces.NewRecord(newColumns)

		// Get the ID
		var id int
		switch v := record.Columns["id"].(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			continue // Skip invalid records
		}

		// Insert into new tree
		newTable.tree.Insert(id, newRecord)
	}

	return newTable
}

// MarshalJSON customizes the JSON representation of the Table
func (t *Table) MarshalJSON() ([]byte, error) {
	// Create a structure to hold both columns and records
	tableData := struct {
		Columns []Column                     `json:"columns"`
		Records map[string]map[string]interface{} `json:"records"`
	}{
		Columns: t.columns,
		Records: make(map[string]map[string]interface{}),
	}

	// Get all records from the B-tree
	records := t.tree.Scan()
	for _, record := range records {
		id := record.Columns["id"]
		tableData.Records[fmt.Sprintf("%v", id)] = record.Columns
	}

	return json.Marshal(tableData)
}

// UnmarshalJSON customizes the JSON deserialization of the Table
func (t *Table) UnmarshalJSON(data []byte) error {
	// Create a temporary structure to hold the data
	var tableData struct {
		Columns []Column                     `json:"columns"`
		Records map[string]map[string]interface{} `json:"records"`
	}

	if err := json.Unmarshal(data, &tableData); err != nil {
		return err
	}

	// Initialize the table
	t.columns = tableData.Columns
	t.tree = NewBTree()

	// Restore records
	for _, recordData := range tableData.Records {
		record := interfaces.NewRecord(recordData)
		if id, ok := recordData["id"].(float64); ok {
			t.tree.Insert(int(id), record)
		}
	}

	return nil
}

package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"sync"

	"sqlight/pkg/interfaces"
)

// Database represents a SQLite database
type Database struct {
	mutex       sync.RWMutex
	tables      map[string]*interfaces.Table
	path        string
	inTransaction bool
	snapshot    map[string]*interfaces.Table
}

// NewDatabase creates a new database instance
func NewDatabase(path string) (*Database, error) {
	db := &Database{
		tables: make(map[string]*interfaces.Table),
		path:   path,
	}

	// Load existing database if file exists
	if _, err := os.Stat(path); err == nil {
		if err := db.load(); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Execute executes a SQL statement
func (d *Database) Execute(stmt interfaces.Statement) (*interfaces.Result, error) {
	switch stmt.(type) {
	case *interfaces.BeginTransactionStatement:
		return d.executeBeginTransaction()
	case *interfaces.CommitStatement:
		return d.executeCommit()
	case *interfaces.RollbackStatement:
		return d.executeRollback()
	default:
		d.mutex.Lock()
		defer d.mutex.Unlock()
		
		switch s := stmt.(type) {
		case *interfaces.CreateStatement:
			return d.executeCreate(s)
		case *interfaces.InsertStatement:
			return d.executeInsert(s)
		case *interfaces.SelectStatement:
			return d.executeSelect(s)
		case *interfaces.DropStatement:
			return d.executeDrop(s)
		case *interfaces.DescribeStatement:
			return d.executeDescribe(s)
		case *interfaces.DeleteStatement:
			return d.executeDelete(s)
		default:
			return nil, fmt.Errorf("unsupported statement type: %T", stmt)
		}
	}
}

// executeBeginTransaction starts a new transaction
func (d *Database) executeBeginTransaction() (*interfaces.Result, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.inTransaction {
		return nil, fmt.Errorf("transaction already in progress")
	}

	// Create a deep copy of current database state
	d.snapshot = make(map[string]*interfaces.Table)
	for name, table := range d.tables {
		newTable := &interfaces.Table{
			Name:    table.Name,
			Columns: make([]interfaces.Column, len(table.Columns)),
			Records: make([]*interfaces.Record, len(table.Records)),
		}
		copy(newTable.Columns, table.Columns)
		for i, record := range table.Records {
			newRecord := &interfaces.Record{
				Columns: make(map[string]interface{}),
			}
			for k, v := range record.Columns {
				newRecord.Columns[k] = v
			}
			newTable.Records[i] = newRecord
		}
		d.snapshot[name] = newTable
	}

	d.inTransaction = true
	return &interfaces.Result{
		Success: true,
		Message: "Transaction started",
	}, nil
}

// executeCommit commits the current transaction
func (d *Database) executeCommit() (*interfaces.Result, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.inTransaction {
		return nil, fmt.Errorf("no transaction in progress")
	}

	// Clear snapshot and commit by saving to disk
	d.tables = d.snapshot
	d.snapshot = nil
	d.inTransaction = false
	if err := d.save(); err != nil {
		return nil, err
	}

	return &interfaces.Result{
		Success: true,
		Message: "Transaction committed successfully",
	}, nil
}

// executeRollback aborts the current transaction
func (d *Database) executeRollback() (*interfaces.Result, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.inTransaction {
		return nil, fmt.Errorf("no transaction in progress")
	}

	// Restore from snapshot
	d.snapshot = nil
	d.inTransaction = false

	return &interfaces.Result{
		Success: true,
		Message: "Transaction rolled back successfully",
	}, nil
}

// executeCreate handles CREATE TABLE statements
func (d *Database) executeCreate(stmt *interfaces.CreateStatement) (*interfaces.Result, error) {
	if _, exists := d.tables[stmt.TableName]; exists {
		return nil, fmt.Errorf("table %s already exists", stmt.TableName)
	}

	// Validate constraints
	primaryKeyCount := 0
	for _, col := range stmt.Columns {
		if col.PrimaryKey {
			primaryKeyCount++
			if primaryKeyCount > 1 {
				return nil, fmt.Errorf("table can only have one PRIMARY KEY")
			}
		}
	}

	// Create table
	table := &interfaces.Table{
		Name:    stmt.TableName,
		Columns: stmt.Columns,
		Records: make([]*interfaces.Record, 0),
	}

	// Add table to transaction if in transaction, otherwise add to database
	if d.inTransaction {
		d.tables[stmt.TableName] = table
	} else {
		d.tables[stmt.TableName] = table
		if err := d.save(); err != nil {
			return nil, err
		}
	}

	return &interfaces.Result{
		Success: true,
		Message: fmt.Sprintf("Table %s created successfully", stmt.TableName),
	}, nil
}

// getColumnValue converts a value to the appropriate type based on column definition
func getColumnValue(colDef *interfaces.Column, value interface{}) (interface{}, error) {
    switch colDef.Type {
    case "INT", "INTEGER":
        switch v := value.(type) {
        case int:
            return v, nil
        case float64:
            return int(v), nil
        case string:
            return strconv.Atoi(v)
        default:
            return nil, fmt.Errorf("invalid integer value: %v", value)
        }
    case "TEXT":
        return fmt.Sprintf("%v", value), nil
    default:
        return value, nil
    }
}

// compareValues compares two values based on their types
func compareValues(v1, v2 interface{}) bool {
    // Handle nil values
    if v1 == nil && v2 == nil {
        return true
    }
    if v1 == nil || v2 == nil {
        return false
    }

    switch val1 := v1.(type) {
    case int:
        switch val2 := v2.(type) {
        case int:
            return val1 == val2
        case float64:
            return float64(val1) == val2
        case string:
            if num, err := strconv.Atoi(val2); err == nil {
                return val1 == num
            }
        }
    case float64:
        switch val2 := v2.(type) {
        case float64:
            return val1 == val2
        case int:
            return val1 == float64(val2)
        case string:
            if num, err := strconv.ParseFloat(val2, 64); err == nil {
                return val1 == num
            }
        }
    case string:
        switch val2 := v2.(type) {
        case string:
            // For string comparison, use exact matching
            return val1 == val2
        case int:
            if num, err := strconv.Atoi(val1); err == nil {
                return num == val2
            }
        case float64:
            if num, err := strconv.ParseFloat(val1, 64); err == nil {
                return num == val2
            }
        }
    }
    return false
}

// executeInsert handles INSERT statements
func (d *Database) executeInsert(stmt *interfaces.InsertStatement) (*interfaces.Result, error) {
    table, tableName, err := d.getTable(stmt.TableName, true)
    if err != nil {
        return nil, err
    }

    // Create column name mapping for case-insensitive comparison
    columnMap := d.getColumnMap(table)

    // Create a new record with the provided values
    record := &interfaces.Record{
        Columns: make(map[string]interface{}),
    }

    // Validate column count
    if len(stmt.Columns) != len(stmt.Values) {
        return nil, fmt.Errorf("column count (%d) does not match value count (%d)", len(stmt.Columns), len(stmt.Values))
    }

    // First pass: validate and set column values
    for i, col := range stmt.Columns {
        actualCol, exists := columnMap[strings.ToLower(col)]
        if !exists {
            return nil, fmt.Errorf("column %s does not exist", col)
        }

        // Get column definition
        var colDef *interfaces.Column
        for _, c := range table.Columns {
            if c.Name == actualCol {
                colDef = &c
                break
            }
        }
        if colDef == nil {
            return nil, fmt.Errorf("column %s not found in table definition", actualCol)
        }

        // Convert and validate value
        value, err := getColumnValue(colDef, stmt.Values[i])
        if err != nil {
            return nil, fmt.Errorf("invalid value for column %s: %v", actualCol, err)
        }
        record.Columns[actualCol] = value
    }

    // Second pass: validate constraints
    for _, col := range table.Columns {
        value, exists := record.Columns[col.Name]

        // Check NOT NULL constraint
        if !col.Nullable && (!exists || value == nil) {
            return nil, fmt.Errorf("column %s cannot be null", col.Name)
        }

        // Check PRIMARY KEY and UNIQUE constraints
        if (col.PrimaryKey || col.Unique) && exists && value != nil {
            for _, existingRecord := range table.Records {
                existingValue := existingRecord.Columns[col.Name]
                if existingValue == nil {
                    continue
                }

                // For string values, do exact comparison
                if strValue, ok := value.(string); ok {
                    if strExistingValue, ok := existingValue.(string); ok {
                        if strValue == strExistingValue {
                            constraint := "UNIQUE"
                            if col.PrimaryKey {
                                constraint = "PRIMARY KEY"
                            }
                            return nil, fmt.Errorf("duplicate value in %s column %s", constraint, col.Name)
                        }
                        continue
                    }
                }

                // For other types use compareValues
                if compareValues(value, existingValue) {
                    constraint := "UNIQUE"
                    if col.PrimaryKey {
                        constraint = "PRIMARY KEY"
                    }
                    return nil, fmt.Errorf("duplicate value in %s column %s", constraint, col.Name)
                }
            }
        }
    }

    // Add record to table
    table.Records = append(table.Records, record)

    // Update the appropriate table map
    if d.inTransaction {
        d.snapshot[tableName] = table
    } else {
        d.tables[tableName] = table
        if err := d.save(); err != nil {
            return nil, err
        }
    }

    return &interfaces.Result{
        Success: true,
        Message: "Record inserted successfully",
    }, nil
}

// executeSelect handles SELECT statements
func (d *Database) executeSelect(stmt *interfaces.SelectStatement) (*interfaces.Result, error) {
    table, _, err := d.getTable(stmt.TableName, true)
    if err != nil {
        return nil, err
    }

    // Get column names case-insensitively
    columnMap := d.getColumnMap(table)
    
    // Prepare result columns
    columns := make([]string, 0)
    if len(stmt.Columns) == 0 || stmt.Columns[0] == "*" {
        for _, col := range table.Columns {
            columns = append(columns, col.Name)
        }
    } else {
        for _, col := range stmt.Columns {
            actualCol, exists := columnMap[strings.ToLower(col)]
            if !exists {
                return nil, fmt.Errorf("column %s does not exist", col)
            }
            columns = append(columns, actualCol)
        }
    }

    // Filter records based on WHERE conditions
    var filteredRecords []*interfaces.Record
    
    // If no WHERE conditions, include all records
    if len(stmt.Where) == 0 {
        filteredRecords = table.Records
    } else {
        // Apply WHERE conditions
        for _, record := range table.Records {
            match := true
            for whereCol, whereCondition := range stmt.Where {
                // Get actual column name from case-insensitive map
                actualCol, exists := columnMap[strings.ToLower(whereCol)]
                if !exists {
                    return nil, fmt.Errorf("column %s does not exist", whereCol)
                }

                recordValue := record.Columns[actualCol]
                if recordValue == nil {
                    match = false
                    break
                }

                // Extract operator and value from the condition
                condMap, ok := whereCondition.(map[string]interface{})
                if !ok {
                    return nil, fmt.Errorf("invalid where condition format")
                }
                
                operator := condMap["operator"].(string)
                whereVal := condMap["value"]

                // Compare based on operator
                if !compareWithOperator(whereVal, recordValue, operator) {
                    match = false
                    break
                }
            }

            if match {
                filteredRecords = append(filteredRecords, record)
            }
        }
    }

    // Format records
    var formattedRecords []*interfaces.Record
    for _, record := range filteredRecords {
        formattedRecord := &interfaces.Record{
            Columns: make(map[string]interface{}),
        }
        for _, col := range columns {
            formattedRecord.Columns[col] = record.Columns[col]
        }
        formattedRecords = append(formattedRecords, formattedRecord)
    }

    return &interfaces.Result{
        Success:  true,
        Columns:  columns,
        Records:  formattedRecords,
        IsSelect: true,
    }, nil
}

// compareWithOperator compares two values using the specified operator
func compareWithOperator(v1, v2 interface{}, operator string) bool {
    // Handle nil values
    if v1 == nil && v2 == nil {
        return operator == "="
    }
    if v1 == nil || v2 == nil {
        return operator == "!="
    }

    // v1 is the value from the WHERE condition
    // v2 is the value from the record
    // So the comparison should be: record_value operator condition_value
    // For example: if WHERE age > 30, then we check if record.age > 30

    switch val1 := v1.(type) {
    case int:
        switch val2 := v2.(type) {
        case int:
            return compareInts(val2, val1, operator)
        case float64:
            return compareFloats(val2, float64(val1), operator)
        case string:
            if num, err := strconv.Atoi(val2); err == nil {
                return compareInts(num, val1, operator)
            }
            if num, err := strconv.ParseFloat(val2, 64); err == nil {
                return compareFloats(num, float64(val1), operator)
            }
        }
    case float64:
        switch val2 := v2.(type) {
        case float64:
            return compareFloats(val2, val1, operator)
        case int:
            return compareFloats(float64(val2), val1, operator)
        case string:
            if num, err := strconv.ParseFloat(val2, 64); err == nil {
                return compareFloats(num, val1, operator)
            }
        }
    case string:
        switch val2 := v2.(type) {
        case string:
            return compareStrings(val2, val1, operator)
        case int:
            if num, err := strconv.Atoi(val1); err == nil {
                return compareInts(val2, num, operator)
            }
        case float64:
            if num, err := strconv.ParseFloat(val1, 64); err == nil {
                return compareFloats(val2, num, operator)
            }
        }
    }
    return false
}

// compareInts compares two integers using the specified operator
func compareInts(a, b int, operator string) bool {
    switch operator {
    case "=":
        return a == b
    case "!=":
        return a != b
    case ">":
        return a > b
    case "<":
        return a < b
    case ">=":
        return a >= b
    case "<=":
        return a <= b
    default:
        return false
    }
}

// compareFloats compares two floats using the specified operator
func compareFloats(a, b float64, operator string) bool {
    switch operator {
    case "=":
        return a == b
    case "!=":
        return a != b
    case ">":
        return a > b
    case "<":
        return a < b
    case ">=":
        return a >= b
    case "<=":
        return a <= b
    default:
        return false
    }
}

// compareStrings compares two strings using the specified operator
func compareStrings(a, b string, operator string) bool {
    switch operator {
    case "=":
        return strings.EqualFold(a, b)
    case "!=":
        return !strings.EqualFold(a, b)
    case ">":
        return a > b
    case "<":
        return a < b
    case ">=":
        return a >= b
    case "<=":
        return a <= b
    default:
        return false
    }
}

// executeDescribe handles DESCRIBE statements
func (d *Database) executeDescribe(stmt *interfaces.DescribeStatement) (*interfaces.Result, error) {
	table, _, err := d.getTable(stmt.TableName, true)
	if err != nil {
		return nil, err
	}

	// Format column information
	columns := []string{"Field", "Type", "Constraints"}
	var records []*interfaces.Record

	for _, col := range table.Columns {
		constraints := make([]string, 0)
		if col.PrimaryKey {
			constraints = append(constraints, "PRIMARY KEY")
		}
		if !col.Nullable {
			constraints = append(constraints, "NOT NULL")
		}
		if col.Unique {
			constraints = append(constraints, "UNIQUE")
		}

		record := &interfaces.Record{
			Columns: map[string]interface{}{
				"Field":       col.Name,
				"Type":        col.Type,
				"Constraints": strings.Join(constraints, ", "),
			},
		}
		records = append(records, record)
	}

	return &interfaces.Result{
		Success:  true,
		Columns:  columns,
		Records:  records,
		IsSelect: true,
	}, nil
}

// executeDrop handles DROP TABLE statements
func (d *Database) executeDrop(stmt *interfaces.DropStatement) (*interfaces.Result, error) {
	_, tableName, err := d.getTable(stmt.TableName, true)
	if err != nil {
		return nil, err
	}

	// Remove table from the appropriate map
	if d.inTransaction {
		delete(d.snapshot, tableName)
	} else {
		delete(d.tables, tableName)
		if err := d.save(); err != nil {
			return nil, err
		}
	}

	return &interfaces.Result{
		Success: true,
		Message: fmt.Sprintf("Table %s dropped successfully", stmt.TableName),
	}, nil
}

// executeDelete handles DELETE statements
func (d *Database) executeDelete(stmt *interfaces.DeleteStatement) (*interfaces.Result, error) {
	table, tableName, err := d.getTable(stmt.TableName, true)
	if err != nil {
		return nil, err
	}

	// Create column name mapping for case-insensitive comparison
	columnMap := d.getColumnMap(table)

	// Filter records that match WHERE conditions
	newRecords := make([]*interfaces.Record, 0)
	deletedCount := 0

	// If no WHERE clause, delete all records
	if len(stmt.Where) == 0 {
		deletedCount = len(table.Records)
		newRecords = make([]*interfaces.Record, 0) // Empty the records
	} else {
		// Process records with WHERE clause
		for _, record := range table.Records {
			match := true
			for whereCol, whereCondition := range stmt.Where {
				// Get actual column name from case-insensitive map
				actualCol, exists := columnMap[strings.ToLower(whereCol)]
				if !exists {
					return nil, fmt.Errorf("column %s does not exist", whereCol)
				}

				recordValue := record.Columns[actualCol]
				if recordValue == nil {
					match = false
					break
				}

                // Extract operator and value from the condition
                condMap, ok := whereCondition.(map[string]interface{})
                if !ok {
                    return nil, fmt.Errorf("invalid where condition format")
                }
                
                operator := condMap["operator"].(string)
                whereVal := condMap["value"]

                // Compare based on operator
                if !compareWithOperator(whereVal, recordValue, operator) {
                    match = false
                    break
                }
			}
			if !match {
				newRecords = append(newRecords, record)
			} else {
				deletedCount++
			}
		}
	}

	// Update table with filtered records
	table.Records = newRecords

	// Update the appropriate table map
	if d.inTransaction {
		d.snapshot[tableName] = table
	} else {
		d.tables[tableName] = table
		if err := d.save(); err != nil {
			return nil, err
		}
	}

	return &interfaces.Result{
		Success: true,
		Message: fmt.Sprintf("%d record(s) deleted successfully", deletedCount),
	}, nil
}

// getTable finds a table case-insensitively
func (d *Database) getTable(tableName string, useSnapshot bool) (*interfaces.Table, string, error) {
	// Get the target table map based on transaction state
	tables := d.tables
	if useSnapshot && d.inTransaction {
		tables = d.snapshot
	}

	// Find table case-insensitively
	var table *interfaces.Table
	var actualName string
	for name, t := range tables {
		if strings.EqualFold(name, tableName) {
			table = t
			actualName = name
			break
		}
	}

	if table == nil {
		return nil, "", fmt.Errorf("table %s does not exist", tableName)
	}

	return table, actualName, nil
}

// getColumnMap creates a case-insensitive column name mapping
func (d *Database) getColumnMap(table *interfaces.Table) map[string]string {
	columnMap := make(map[string]string)
	for _, col := range table.Columns {
		columnMap[strings.ToLower(col.Name)] = col.Name
	}
	return columnMap
}

// save saves the database to a file
func (d *Database) save() error {
	data, err := json.MarshalIndent(d.tables, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.path, data, 0644)
}

// load loads the database from a file
func (d *Database) load() error {
	data, err := ioutil.ReadFile(d.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &d.tables)
}

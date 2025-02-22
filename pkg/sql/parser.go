package sql

import (
	"fmt"
	"os"
	"sqlite-clone/pkg/interfaces"
	"strconv"
	"strings"
)

// Update Statement to use the interface
type Statement interface {
	Exec(db interfaces.Database) error
}

// CreateStatement represents a CREATE TABLE command
type CreateStatement struct {
	Table   string
	Columns []interfaces.ColumnDef
}

func (s *CreateStatement) Exec(db interfaces.Database) error {
	return db.CreateTable(s.Table, s.Columns)
}

// InsertStatement represents an INSERT command
type InsertStatement struct {
	Table  string
	Values []interface{}
}

func (s *InsertStatement) Exec(db interfaces.Database) error {
	// Get column information
	columns, err := db.GetTableColumns(s.Table)
	if err != nil {
		return err
	}

	// Create a map for the record's columns
	recordColumns := make(map[string]interface{})
	for i, col := range columns {
		if i < len(s.Values) {
			recordColumns[col] = s.Values[i]
		}
	}

	// Create new record with the columns map
	record := interfaces.NewRecord(recordColumns)

	return db.InsertIntoTable(s.Table, record)
}

// SelectStatement represents a SELECT query
type SelectStatement struct {
	Table      string
	WhereCol   string
	WhereValue interface{}
}

func (s *SelectStatement) Exec(db interfaces.Database) error {
	// Get all records
	records, err := db.SelectFromTable(s.Table)
	if err != nil {
		return err
	}

	// Get column names
	columns, err := db.GetTableColumns(s.Table)
	if err != nil {
		return err
	}

	// Filter records if where clause exists
	var filteredRecords []interface{}
	if s.WhereCol != "" {
		for _, record := range records {
			if r, ok := record.(*interfaces.Record); ok {
				if val, exists := r.Columns[s.WhereCol]; exists {
					// Convert both values to strings for comparison
					valStr := fmt.Sprintf("%v", val)
					whereStr := fmt.Sprintf("%v", s.WhereValue)
					if valStr == whereStr {
						filteredRecords = append(filteredRecords, record)
					}
				}
			}
		}
	} else {
		filteredRecords = records
	}

	// Print records
	if len(filteredRecords) > 0 {
		printRecords(columns, filteredRecords)
	}
	return nil
}

// ExitStatement represents an EXIT command
type ExitStatement struct{}

func (s *ExitStatement) Exec(db interfaces.Database) error {
	fmt.Println("Exiting the SQLite Clone. Goodbye!")
	os.Exit(0)
	return nil
}

// UpdateStatement represents an UPDATE command
type UpdateStatement struct {
	TableName   string
	SetColumns  map[string]interface{}
	WhereColumn string
	WhereValue  interface{}
}

func (s *UpdateStatement) Exec(db interfaces.Database) error {
	return db.UpdateTable(s.TableName, s.SetColumns, s.WhereColumn, s.WhereValue)
}

// DeleteStatement represents a DELETE command
type DeleteStatement struct {
	TableName string
	WhereColumn string
	WhereValue interface{}
}

func (s *DeleteStatement) Exec(db interfaces.Database) error {
	return db.DeleteFromTable(s.TableName, s.WhereColumn, s.WhereValue)
}

// Helper function to print records in table format
func printRecords(columns []string, records []interface{}) {
	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(col)
	}

	// Find maximum width for each column
	for _, record := range records {
		if r, ok := record.(*interfaces.Record); ok {
			for col, value := range r.Columns {
				width := len(fmt.Sprintf("%v", value))
				if width > widths[col] {
					widths[col] = width
				}
			}
		}
	}

	// Print header
	printSeparator(columns, widths)
	printRow(columns, columns, widths) // Column names as header
	printSeparator(columns, widths)

	// Print records
	for _, record := range records {
		if r, ok := record.(*interfaces.Record); ok {
			values := make([]string, len(columns))
			for i, col := range columns {
				values[i] = fmt.Sprintf("%v", r.Columns[col])
			}
			printRow(columns, values, widths)
		}
	}
	printSeparator(columns, widths)
}

func printRecord(columns []string, record interface{}) {
	printRecords(columns, []interface{}{record})
}

func printSeparator(columns []string, widths map[string]int) {
	for _, col := range columns {
		fmt.Print("+")
		fmt.Print(strings.Repeat("-", widths[col]+2))
	}
	fmt.Println("+")
}

func printRow(columns []string, values []string, widths map[string]int) {
	for i, col := range columns {
		fmt.Printf("| %-*s ", widths[col], values[i])
	}
	fmt.Println("|")
}

// ParseSQL parses SQL commands
func ParseSQL(query string) (interfaces.Statement, error) {
	// Convert to uppercase for case-insensitive matching of keywords only
	upperQuery := strings.ToUpper(query)
	if strings.HasPrefix(upperQuery, "CREATE TABLE") {
		return parseCreateTable(query)
	} else if strings.HasPrefix(upperQuery, "INSERT INTO") {
		return parseInsert(query)
	} else if strings.HasPrefix(upperQuery, "SELECT") {
		return parseSelect(query)
	} else if strings.HasPrefix(upperQuery, "DELETE FROM") {
		return parseDelete(query)
	} else if strings.HasPrefix(upperQuery, "UPDATE") {
		return parseUpdate(query)
	}
	return nil, fmt.Errorf("unsupported SQL command")
}

func parseCreateTable(query string) (*CreateStatement, error) {
	// Parse CREATE TABLE query
	words := strings.Fields(query)
	if len(words) < 4 || strings.ToUpper(words[1]) != "TABLE" {
		return nil, fmt.Errorf("invalid CREATE TABLE syntax")
	}

	tableName := words[2]
	columnsStr := strings.Join(words[3:], " ")

	// Parse column definitions
	if !strings.HasPrefix(columnsStr, "(") || !strings.HasSuffix(columnsStr, ")") {
		return nil, fmt.Errorf("column definitions must be enclosed in parentheses")
	}

	// Extract column names and types
	columnsStr = columnsStr[1 : len(columnsStr)-1]
	columnDefs := strings.Split(columnsStr, ",")
	columns := make([]interfaces.ColumnDef, len(columnDefs))
	for i, def := range columnDefs {
		parts := strings.Fields(def)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid column definition: %s", def)
		}
		columns[i] = interfaces.ColumnDef{
			Name: parts[0],
			Type: strings.ToUpper(parts[1]), // Store the type in uppercase
		}
	}

	return &CreateStatement{Table: tableName, Columns: columns}, nil
}

func parseInsert(query string) (*InsertStatement, error) {
	// Remove trailing semicolon if present
	query = strings.TrimSuffix(query, ";")

	// Parse INSERT query
	words := strings.Fields(query)
	if len(words) < 4 || strings.ToUpper(words[1]) != "INTO" || strings.ToUpper(words[3]) != "VALUES" {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	tableName := strings.TrimSpace(words[2])
	tableName = strings.ToLower(tableName)

	// Join the remaining words for values parsing
	valuesStr := strings.Join(words[4:], " ")
	valuesStr = strings.TrimSpace(valuesStr)

	// Handle both formats: with or without parentheses
	if !strings.HasPrefix(valuesStr, "(") {
		valuesStr = "(" + valuesStr + ")"
	}

	// Parse values
	if !strings.HasPrefix(valuesStr, "(") || !strings.HasSuffix(valuesStr, ")") {
		return nil, fmt.Errorf("values must be enclosed in parentheses")
	}

	// Extract values between parentheses
	valuesStr = valuesStr[1 : len(valuesStr)-1]
	valueParts := strings.Split(valuesStr, ",")
	values := make([]interface{}, len(valueParts))

	for i, part := range valueParts {
		part = strings.TrimSpace(part)
		// Handle string values (quoted)
		if strings.HasPrefix(part, "'") && strings.HasSuffix(part, "'") {
			values[i] = strings.Trim(part, "'")
		} else {
			// Try to parse as float64 for numeric values
			if floatVal, err := strconv.ParseFloat(part, 64); err == nil {
				values[i] = floatVal
			} else {
				return nil, fmt.Errorf("invalid value: %s", part)
			}
		}
	}

	return &InsertStatement{
		Table:  tableName,
		Values: values,
	}, nil
}

func parseSelect(query string) (*SelectStatement, error) {
	// Remove trailing semicolon if present
	query = strings.TrimSuffix(query, ";")

	// Parse SELECT query
	words := strings.Fields(query)
	if len(words) < 3 || strings.ToUpper(words[0]) != "SELECT" {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}

	// Check for * and FROM
	if words[1] != "*" || strings.ToUpper(words[2]) != "FROM" {
		return nil, fmt.Errorf("only SELECT * FROM is supported")
	}

	if len(words) < 4 {
		return nil, fmt.Errorf("table name required after FROM")
	}

	tableName := words[3]
	stmt := &SelectStatement{Table: tableName}

	// Check for WHERE clause
	if len(words) > 4 {
		if len(words) < 8 || strings.ToUpper(words[4]) != "WHERE" || words[6] != "=" {
			return nil, fmt.Errorf("WHERE clause must be in format: WHERE column = value")
		}

		whereCol := strings.ToLower(words[5])
		whereVal := words[7]

		// Handle quoted string values
		if strings.HasPrefix(whereVal, "'") && strings.HasSuffix(whereVal, "'") {
			stmt.WhereCol = whereCol
			stmt.WhereValue = whereVal[1 : len(whereVal)-1]
		} else {
			// Try to convert to number if not a string
			if num, err := strconv.ParseFloat(whereVal, 64); err == nil {
				stmt.WhereCol = whereCol
				stmt.WhereValue = num
			} else {
				stmt.WhereCol = whereCol
				stmt.WhereValue = whereVal
			}
		}
	}

	return stmt, nil
}

func parseDelete(query string) (*DeleteStatement, error) {
	// Remove trailing semicolon if present
	query = strings.TrimSuffix(query, ";")

	// Parse DELETE query
	words := strings.Fields(query)
	if len(words) < 3 || strings.ToUpper(words[1]) != "FROM" {
		return nil, fmt.Errorf("invalid DELETE syntax")
	}

	tableName := strings.ToLower(words[2])
	stmt := &DeleteStatement{TableName: tableName}

	// Parse WHERE clause if present
	if len(words) > 3 {
		if strings.ToUpper(words[3]) != "WHERE" || len(words) < 7 {
			return nil, fmt.Errorf("invalid WHERE clause in DELETE")
		}
		stmt.WhereColumn = strings.ToLower(words[4])
		if words[5] != "=" {
			return nil, fmt.Errorf("WHERE clause must use = operator")
		}
		value := words[6]

		// Handle quoted string values
		if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			stmt.WhereValue = value[1:len(value)-1]
		} else {
			// Try to convert to number if not a string
			if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				stmt.WhereValue = floatVal
			} else {
				stmt.WhereValue = value
			}
		}
	}

	return stmt, nil
}

func parseUpdate(query string) (*UpdateStatement, error) {
	// Remove trailing semicolon if present
	query = strings.TrimSuffix(query, ";")

	// UPDATE table_name SET column1 = value1, column2 = value2 WHERE column = value
	parts := strings.Split(query, "WHERE")
	if len(parts) != 2 {
		return nil, fmt.Errorf("UPDATE statement must have a WHERE clause")
	}

	// Parse the main part (before WHERE)
	mainParts := strings.Split(parts[0], "SET")
	if len(mainParts) != 2 {
		return nil, fmt.Errorf("invalid UPDATE statement format")
	}

	// Get table name
	tablePart := strings.TrimSpace(mainParts[0])
	if !strings.HasPrefix(strings.ToUpper(tablePart), "UPDATE") {
		return nil, fmt.Errorf("invalid UPDATE statement format")
	}
	tableName := strings.TrimSpace(strings.TrimPrefix(tablePart, "UPDATE"))
	tableName = strings.ToLower(tableName) // Convert table name to lowercase

	// Parse SET clause
	setColumns := make(map[string]interface{})
	setPairs := strings.Split(strings.TrimSpace(mainParts[1]), ",")
	for _, pair := range setPairs {
		keyVal := strings.Split(strings.TrimSpace(pair), "=")
		if len(keyVal) != 2 {
			return nil, fmt.Errorf("invalid SET clause format")
		}
		key := strings.TrimSpace(keyVal[0])
		key = strings.ToLower(key) // Convert column name to lowercase
		val := strings.TrimSpace(keyVal[1])

		// Try to parse as float64 for numeric values
		if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
			setColumns[key] = floatVal
		} else {
			// Remove quotes for string values
			setColumns[key] = strings.Trim(val, "'\"")
		}
	}

	// Parse WHERE clause
	whereClause := strings.TrimSpace(parts[1])
	whereParts := strings.Split(whereClause, "=")
	if len(whereParts) != 2 {
		return nil, fmt.Errorf("invalid WHERE clause in UPDATE statement")
	}

	whereColumn := strings.TrimSpace(whereParts[0])
	whereColumn = strings.ToLower(whereColumn) // Convert where column to lowercase
	whereVal := strings.TrimSpace(whereParts[1])

	stmt := &UpdateStatement{
		TableName:   tableName,
		SetColumns:  setColumns,
		WhereColumn: whereColumn,
	}

	// Try to parse where value as float64 for numeric values
	if floatVal, err := strconv.ParseFloat(whereVal, 64); err == nil {
		stmt.WhereValue = floatVal
	} else {
		// Remove quotes and any trailing semicolon for string values
		whereVal = strings.TrimSuffix(whereVal, ";")
		stmt.WhereValue = strings.Trim(whereVal, "'\"")
	}

	return stmt, nil
}

// Helper function to parse values
func parseValues(input string) ([]interface{}, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "(") || !strings.HasSuffix(input, ")") {
		return nil, fmt.Errorf("values must be enclosed in parentheses")
	}

	// Remove parentheses
	input = input[1 : len(input)-1]
	parts := strings.Split(input, ",")

	values := make([]interface{}, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		// Handle string values (quoted)
		if strings.HasPrefix(part, "'") && strings.HasSuffix(part, "'") {
			values[i] = part[1 : len(part)-1]
		} else {
			// Number value
			if num, err := strconv.Atoi(part); err == nil {
				values[i] = num
			} else {
				return nil, fmt.Errorf("invalid number: %s", part)
			}
		}
	}

	return values, nil
}

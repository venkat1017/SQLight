package sql

import (
	"fmt"
	"os"
	"sqlight/pkg/interfaces"
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
	records, err := db.SelectFromTable(s.Table, s.WhereCol, s.WhereValue)
	if err != nil {
		return err
	}

	// If no records found
	if len(records) == 0 {
		fmt.Println("No records found")
		return nil
	}

	// Get column information
	columns, err := db.GetTableColumns(s.Table)
	if err != nil {
		return err
	}

	// Print records
	printRecords(columns, records)
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
	TableName   string
	WhereColumn string
	WhereValue  interface{}
}

func (s *DeleteStatement) Exec(db interfaces.Database) error {
	return db.DeleteFromTable(s.TableName, s.WhereColumn, s.WhereValue)
}

// BeginStatement represents a BEGIN TRANSACTION command
type BeginStatement struct{}

func (s *BeginStatement) Exec(db interfaces.Database) error {
	return db.Begin()
}

// CommitStatement represents a COMMIT command
type CommitStatement struct{}

func (s *CommitStatement) Exec(db interfaces.Database) error {
	return db.Commit()
}

// RollbackStatement represents a ROLLBACK command
type RollbackStatement struct{}

func (s *RollbackStatement) Exec(db interfaces.Database) error {
	return db.Rollback()
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

// Parser represents a SQL parser
type Parser struct {
	db interfaces.Database
}

// NewParser creates a new SQL parser
func NewParser(db interfaces.Database) *Parser {
	return &Parser{db: db}
}

// Parse parses and executes a SQL command
func (p *Parser) Parse(query string) (string, error) {
	stmt, err := ParseSQL(query)
	if err != nil {
		return "", err
	}

	err = stmt.Exec(p.db)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

// ParseSQL parses SQL commands
func ParseSQL(query string) (Statement, error) {
	// Remove semicolon if present
	query = strings.TrimSuffix(query, ";")
	query = strings.TrimSpace(query)

	// Convert to uppercase for case-insensitive comparison
	upperQuery := strings.ToUpper(query)

	if strings.HasPrefix(upperQuery, "CREATE TABLE") {
		return parseCreateTable(query)
	} else if strings.HasPrefix(upperQuery, "INSERT INTO") {
		return parseInsert(query)
	} else if strings.HasPrefix(upperQuery, "SELECT") {
		return parseSelect(query)
	} else if strings.HasPrefix(upperQuery, "UPDATE") {
		return parseUpdate(query)
	} else if strings.HasPrefix(upperQuery, "DELETE FROM") {
		return parseDelete(query)
	} else if upperQuery == "BEGIN TRANSACTION" || upperQuery == "BEGIN" {
		return &BeginStatement{}, nil
	} else if upperQuery == "COMMIT" {
		return &CommitStatement{}, nil
	} else if upperQuery == "ROLLBACK" {
		return &RollbackStatement{}, nil
	} else if upperQuery == "EXIT" {
		return &ExitStatement{}, nil
	}

	return nil, fmt.Errorf("unknown command: %s", query)
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
	if len(words) < 4 || strings.ToUpper(words[1]) != "INTO" {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	tableName := words[2]

	// Find the VALUES keyword
	valuesIndex := -1
	for i, word := range words {
		if strings.ToUpper(word) == "VALUES" {
			valuesIndex = i
			break
		}
	}

	if valuesIndex == -1 || valuesIndex+1 >= len(words) {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	// Parse the values
	valuesStr := strings.Join(words[valuesIndex+1:], " ")
	parsedValues, err := parseValues(valuesStr)
	if err != nil {
		return nil, err
	}

	return &InsertStatement{Table: tableName, Values: parsedValues}, nil
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
			stmt.WhereValue = value[1 : len(value)-1]
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

	// UPDATE table_name SET column1 = value1 WHERE column = value
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

	// Parse SET clause
	setColumns := make(map[string]interface{})
	setPairs := strings.Split(strings.TrimSpace(mainParts[1]), ",")
	for _, pair := range setPairs {
		keyVal := strings.Split(strings.TrimSpace(pair), "=")
		if len(keyVal) != 2 {
			return nil, fmt.Errorf("invalid SET clause format")
		}
		key := strings.TrimSpace(keyVal[0])
		val := strings.TrimSpace(keyVal[1])

		// Handle quoted string values
		if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
			setColumns[key] = val[1 : len(val)-1]
		} else {
			// Try to parse as float64 for numeric values
			if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
				setColumns[key] = floatVal
			} else {
				setColumns[key] = val
			}
		}
	}

	// Parse WHERE clause
	whereClause := strings.TrimSpace(parts[1])
	whereParts := strings.Split(whereClause, "=")
	if len(whereParts) != 2 {
		return nil, fmt.Errorf("invalid WHERE clause format")
	}

	whereColumn := strings.TrimSpace(whereParts[0])
	whereValue := strings.TrimSpace(whereParts[1])

	// Handle quoted string values for whereValue
	var whereValueInterface interface{}
	if strings.HasPrefix(whereValue, "'") && strings.HasSuffix(whereValue, "'") {
		whereValueInterface = whereValue[1 : len(whereValue)-1]
	} else {
		// Try to parse as float64 for numeric values
		if floatVal, err := strconv.ParseFloat(whereValue, 64); err == nil {
			whereValueInterface = floatVal
		} else {
			whereValueInterface = whereValue
		}
	}

	return &UpdateStatement{
		TableName:   tableName,
		SetColumns:  setColumns,
		WhereColumn: whereColumn,
		WhereValue:  whereValueInterface,
	}, nil
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

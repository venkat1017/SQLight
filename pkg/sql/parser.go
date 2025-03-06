package sql

import (
	"fmt"
	"regexp"
	"sqlight/pkg/interfaces"
	"strings"
)

// ParseSQL parses a SQL statement and returns a Statement interface
func ParseSQL(sql string) (interfaces.Statement, error) {
	// Trim whitespace and semicolon
	sql = strings.TrimSpace(sql)
	sql = strings.TrimSuffix(sql, ";")

	// Convert to uppercase for easier parsing
	upperSQL := strings.ToUpper(sql)

	if strings.HasPrefix(upperSQL, "CREATE TABLE") {
		return parseCreateTable(sql)
	} else if strings.HasPrefix(upperSQL, "INSERT INTO") {
		return parseInsert(sql)
	} else if strings.HasPrefix(upperSQL, "SELECT") {
		return parseSelect(sql)
	}

	return nil, fmt.Errorf("unsupported SQL statement")
}

// parseCreateTable parses a CREATE TABLE statement
func parseCreateTable(sql string) (*interfaces.CreateStatement, error) {
	// Regular expression for CREATE TABLE
	re := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\((.*)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid CREATE TABLE syntax")
	}

	tableName := matches[1]
	columnDefs := matches[2]

	// Split column definitions
	columns := strings.Split(columnDefs, ",")
	if len(columns) == 0 {
		return nil, fmt.Errorf("no columns specified")
	}

	// Parse each column
	parsedColumns := make([]interfaces.ColumnDef, 0, len(columns))
	for _, col := range columns {
		// Split column definition into parts
		parts := strings.Fields(strings.TrimSpace(col))
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid column definition: %s", col)
		}

		// Create column definition
		column := interfaces.ColumnDef{
			Name: parts[0],
			Type: strings.ToUpper(parts[1]),
		}

		// Parse constraints
		for i := 2; i < len(parts); i++ {
			constraint := strings.ToUpper(parts[i])
			switch constraint {
			case "PRIMARY":
				if i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "KEY" {
					column.PrimaryKey = true
					i++ // Skip next part
				}
			case "NOT":
				if i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "NULL" {
					column.NotNull = true
					i++ // Skip next part
				}
			case "UNIQUE":
				column.Unique = true
			case "REFERENCES":
				if i+1 < len(parts) {
					column.References = parts[i+1]
					i++ // Skip next part
				}
			}
		}

		parsedColumns = append(parsedColumns, column)
	}

	return &interfaces.CreateStatement{
		TableName: tableName,
		Columns:   parsedColumns,
	}, nil
}

// parseInsert parses an INSERT statement
func parseInsert(sql string) (*interfaces.InsertStatement, error) {
	// Regular expression for INSERT
	re := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)\s*\((.*?)\)\s*VALUES\s*\((.*?)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	tableName := matches[1]
	columns := strings.Split(matches[2], ",")
	values := strings.Split(matches[3], ",")

	if len(columns) != len(values) {
		return nil, fmt.Errorf("number of columns does not match number of values")
	}

	// Create values map
	valueMap := make(map[string]interface{})
	for i := range columns {
		col := strings.TrimSpace(columns[i])
		val := strings.TrimSpace(values[i])

		// Remove quotes from string values
		if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
			val = strings.Trim(val, "'")
			valueMap[col] = val
		} else if val == "NULL" {
			valueMap[col] = nil
		} else {
			// Try to parse as integer
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
				valueMap[col] = intVal
			} else {
				valueMap[col] = val
			}
		}
	}

	return &interfaces.InsertStatement{
		TableName: tableName,
		Values:    valueMap,
	}, nil
}

// parseSelect parses a SELECT statement
func parseSelect(sql string) (*interfaces.SelectStatement, error) {
	// Regular expression for SELECT
	re := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM\s+(\w+)(?:\s+WHERE\s+(.*))?`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}

	// Parse columns
	var columns []string
	if matches[1] == "*" {
		columns = []string{"*"}
	} else {
		columns = strings.Split(matches[1], ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}
	}

	// Create statement
	stmt := &interfaces.SelectStatement{
		TableName: matches[2],
		Columns:   columns,
	}

	// Add WHERE clause if present
	if len(matches) > 3 && matches[3] != "" {
		stmt.Where = matches[3]
	}

	return stmt, nil
}

package sql

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "sqlight/pkg/interfaces"
)

// Parse parses a SQL statement and returns the corresponding Statement interface
func Parse(sql string) (interfaces.Statement, error) {
    // Trim whitespace and remove comments
    sql = removeComments(strings.TrimSpace(sql))
    if sql == "" {
        return nil, nil
    }

    // Convert to uppercase for command matching
    upperSQL := strings.ToUpper(sql)

    if strings.HasPrefix(upperSQL, "CREATE TABLE") {
        return parseCreateTable(sql)
    } else if strings.HasPrefix(upperSQL, "INSERT INTO") {
        return parseInsert(sql)
    } else if strings.HasPrefix(upperSQL, "SELECT") {
        return parseSelect(sql)
    } else if strings.HasPrefix(upperSQL, "DROP TABLE") {
        return parseDrop(sql)
    } else if strings.HasPrefix(upperSQL, "DESCRIBE") {
        return parseDescribe(sql)
    } else if strings.HasPrefix(upperSQL, "DELETE FROM") {
        return parseDelete(sql)
    } else if strings.HasPrefix(upperSQL, "BEGIN TRANSACTION") || strings.HasPrefix(upperSQL, "BEGIN") {
        return &interfaces.BeginTransactionStatement{}, nil
    } else if strings.HasPrefix(upperSQL, "COMMIT") {
        return &interfaces.CommitStatement{}, nil
    } else if strings.HasPrefix(upperSQL, "ROLLBACK") {
        return &interfaces.RollbackStatement{}, nil
    }

    return nil, fmt.Errorf("unsupported SQL statement")
}

// removeComments removes SQL comments from the input
func removeComments(sql string) string {
    lines := strings.Split(sql, "\n")
    var result []string
    for _, line := range lines {
        // Remove inline comments
        if idx := strings.Index(line, "--"); idx >= 0 {
            line = strings.TrimSpace(line[:idx])
        }
        if line != "" {
            result = append(result, line)
        }
    }
    return strings.Join(result, "\n")
}

func parseCreateTable(sql string) (*interfaces.CreateStatement, error) {
    // Replace newlines with spaces to handle multi-line statements
    sql = strings.ReplaceAll(sql, "\n", " ")
    
    re := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\((.*)\)`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) != 3 {
        return nil, fmt.Errorf("invalid CREATE TABLE syntax")
    }

    tableName := matches[1]
    columnDefs := strings.Split(matches[2], ",")
    columns := make([]interfaces.Column, 0)

    for _, colDef := range columnDefs {
        colDef = strings.TrimSpace(colDef)
        parts := strings.Fields(colDef)
        if len(parts) < 2 {
            return nil, fmt.Errorf("invalid column definition: %s", colDef)
        }

        col := interfaces.Column{
            Name:     parts[0],
            Type:     strings.ToUpper(parts[1]),
            Nullable: true,
        }

        // Parse constraints
        for i := 2; i < len(parts); i++ {
            constraint := strings.ToUpper(parts[i])
            switch constraint {
            case "PRIMARY":
                if i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "KEY" {
                    col.PrimaryKey = true
                    i++
                }
            case "NOT":
                if i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "NULL" {
                    col.Nullable = false
                    i++
                }
            case "UNIQUE":
                col.Unique = true
            }
        }

        columns = append(columns, col)
    }

    return &interfaces.CreateStatement{
        TableName: tableName,
        Columns:   columns,
    }, nil
}

func parseInsert(sql string) (*interfaces.InsertStatement, error) {
    re := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)\s*\((.*?)\)\s*VALUES\s*\((.*?)\)`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) != 4 {
        return nil, fmt.Errorf("invalid INSERT syntax")
    }

    tableName := matches[1]
    columnStr := matches[2]
    valueStr := matches[3]

    columns := make([]string, 0)
    for _, col := range strings.Split(columnStr, ",") {
        columns = append(columns, strings.TrimSpace(col))
    }

    values := make([]interface{}, 0)
    for _, val := range strings.Split(valueStr, ",") {
        val = strings.TrimSpace(val)
        
        // Handle string values
        if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
            values = append(values, strings.Trim(val, "'"))
            continue
        }
        
        // Handle numeric values
        if num, err := strconv.Atoi(val); err == nil {
            values = append(values, num)
            continue
        }
        
        // Default to string value
        values = append(values, val)
    }

    return &interfaces.InsertStatement{
        TableName: tableName,
        Columns:   columns,
        Values:    values,
    }, nil
}

func parseSelect(sql string) (*interfaces.SelectStatement, error) {
    // Remove trailing semicolon if present
    sql = strings.TrimSuffix(sql, ";")
    
    // Parse table name and columns
    re := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM\s+(\w+)(?:\s+WHERE\s+(.*))?`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) < 3 {
        return nil, fmt.Errorf("invalid SELECT statement syntax")
    }

    // Parse columns
    columns := make([]string, 0)
    for _, col := range strings.Split(matches[1], ",") {
        columns = append(columns, strings.TrimSpace(col))
    }

    // Parse WHERE conditions
    where := make(map[string]interface{})
    if len(matches) > 3 && matches[3] != "" {
        wherePart := strings.TrimSpace(matches[3])
        
        // Split conditions by AND if present
        whereConditions := strings.Split(wherePart, " AND ")
        for _, condition := range whereConditions {
            condition = strings.TrimSpace(condition)
            
            // Check for different comparison operators: =, >, <, >=, <=, !=
            var operator string
            var parts []string
            
            if strings.Contains(condition, ">=") {
                parts = strings.Split(condition, ">=")
                operator = ">="
            } else if strings.Contains(condition, "<=") {
                parts = strings.Split(condition, "<=")
                operator = "<="
            } else if strings.Contains(condition, "!=") {
                parts = strings.Split(condition, "!=")
                operator = "!="
            } else if strings.Contains(condition, ">") {
                parts = strings.Split(condition, ">")
                operator = ">"
            } else if strings.Contains(condition, "<") {
                parts = strings.Split(condition, "<")
                operator = "<"
            } else if strings.Contains(condition, "=") {
                parts = strings.Split(condition, "=")
                operator = "="
            } else {
                return nil, fmt.Errorf("invalid WHERE condition: %s", condition)
            }
            
            if len(parts) != 2 {
                return nil, fmt.Errorf("invalid WHERE condition: %s", condition)
            }

            col := strings.TrimSpace(parts[0])
            val := strings.TrimSpace(parts[1])
            
            // Handle quoted string values (both single and double quotes)
            if (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) ||
               (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) {
                where[col] = map[string]interface{}{
                    "operator": operator,
                    "value": strings.Trim(val, "'\""),
                }
                continue
            }
            
            // Handle numeric values
            if num, err := strconv.Atoi(val); err == nil {
                where[col] = map[string]interface{}{
                    "operator": operator,
                    "value": num,
                }
                continue
            }
            
            // Try to parse as float if not an integer
            if num, err := strconv.ParseFloat(val, 64); err == nil {
                where[col] = map[string]interface{}{
                    "operator": operator,
                    "value": num,
                }
                continue
            }
            
            // Default to string value without quotes
            where[col] = map[string]interface{}{
                "operator": operator,
                "value": val,
            }
        }
    }

    return &interfaces.SelectStatement{
        TableName: matches[2],
        Columns:   columns,
        Where:     where,
    }, nil
}

func parseDrop(sql string) (*interfaces.DropStatement, error) {
    re := regexp.MustCompile(`(?i)DROP\s+TABLE\s+(\w+)`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) != 2 {
        return nil, fmt.Errorf("invalid DROP TABLE syntax")
    }

    return &interfaces.DropStatement{
        TableName: matches[1],
    }, nil
}

func parseDescribe(sql string) (*interfaces.DescribeStatement, error) {
    re := regexp.MustCompile(`(?i)DESCRIBE\s+(\w+)`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) != 2 {
        return nil, fmt.Errorf("invalid DESCRIBE syntax")
    }

    return &interfaces.DescribeStatement{
        TableName: matches[1],
    }, nil
}

func parseDelete(sql string) (*interfaces.DeleteStatement, error) {
    // Remove trailing semicolon if present
    sql = strings.TrimSuffix(sql, ";")
    
    // Parse table name
    re := regexp.MustCompile(`(?i)DELETE\s+FROM\s+(\w+)(?:\s+WHERE\s+(.*))?`)
    matches := re.FindStringSubmatch(sql)
    if len(matches) < 2 {
        return nil, fmt.Errorf("invalid DELETE statement syntax")
    }

    tableName := matches[1]
    conditions := make(map[string]interface{})

    // Parse WHERE conditions if present
    if len(matches) > 2 && matches[2] != "" {
        wherePart := strings.TrimSpace(matches[2])
        
        // Split conditions by AND if present
        whereConditions := strings.Split(wherePart, " AND ")
        for _, condition := range whereConditions {
            condition = strings.TrimSpace(condition)
            
            // Check for different comparison operators: =, >, <, >=, <=, !=
            var operator string
            var parts []string
            
            if strings.Contains(condition, ">=") {
                parts = strings.Split(condition, ">=")
                operator = ">="
            } else if strings.Contains(condition, "<=") {
                parts = strings.Split(condition, "<=")
                operator = "<="
            } else if strings.Contains(condition, "!=") {
                parts = strings.Split(condition, "!=")
                operator = "!="
            } else if strings.Contains(condition, ">") {
                parts = strings.Split(condition, ">")
                operator = ">"
            } else if strings.Contains(condition, "<") {
                parts = strings.Split(condition, "<")
                operator = "<"
            } else if strings.Contains(condition, "=") {
                parts = strings.Split(condition, "=")
                operator = "="
            } else {
                return nil, fmt.Errorf("invalid WHERE condition: %s", condition)
            }
            
            if len(parts) != 2 {
                return nil, fmt.Errorf("invalid WHERE condition: %s", condition)
            }

            column := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])

            // Handle quoted string values (both single and double quotes)
            if (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) ||
               (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) {
                conditions[column] = map[string]interface{}{
                    "operator": operator,
                    "value": strings.Trim(value, "'\""),
                }
                continue
            }
            
            // Handle numeric values
            if num, err := strconv.Atoi(value); err == nil {
                conditions[column] = map[string]interface{}{
                    "operator": operator,
                    "value": num,
                }
                continue
            }
            
            // Try to parse as float if not an integer
            if num, err := strconv.ParseFloat(value, 64); err == nil {
                conditions[column] = map[string]interface{}{
                    "operator": operator,
                    "value": num,
                }
                continue
            }
            
            // Default to string value without quotes
            conditions[column] = map[string]interface{}{
                "operator": operator,
                "value": value,
            }
        }
    }

    return &interfaces.DeleteStatement{
        TableName: tableName,
        Where:     conditions,
    }, nil
}

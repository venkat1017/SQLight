package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "regexp"
    "strings"

    "sqlight/pkg/db"
    "sqlight/pkg/sql"
)

func main() {
    // Print welcome message
    printWelcome()

    // Initialize database
    database, err := db.NewDatabase("database.json")
    if err != nil {
        fmt.Printf("Error initializing database: %v\n", err)
        return
    }

    // Check if a SQL file was provided as an argument
    if len(os.Args) > 1 {
        sqlFile := os.Args[1]
        fmt.Printf("Executing SQL file: %s\n\n", sqlFile)
        
        // Read the file
        content, err := ioutil.ReadFile(sqlFile)
        if err != nil {
            fmt.Printf("Error reading SQL file: %v\n", err)
            return
        }
        
        // Process the file content
        fileContent := string(content)
        
        // Remove comments
        re := regexp.MustCompile(`--.*`)
        fileContent = re.ReplaceAllString(fileContent, "")
        
        // Split into statements
        re = regexp.MustCompile(`;[\s\n]*`)
        statements := re.Split(fileContent, -1)
        
        // Execute each statement
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" {
                continue
            }
            
            // Add semicolon back for parsing
            stmt += ";"
            
            fmt.Printf("Executing: %s\n", stmt)
            
            // Parse and execute
            parsedStmt, err := sql.Parse(stmt)
            if err != nil {
                fmt.Printf("Error parsing statement: %v\n", err)
                continue
            }
            
            // Skip empty statements (comments)
            if parsedStmt == nil {
                continue
            }
            
            // Execute statement
            result, err := database.Execute(parsedStmt)
            if err != nil {
                fmt.Printf("Error executing statement: %v\n", err)
                continue
            }
            
            // Print result
            if result.IsSelect {
                // Print table header
                fmt.Print("| ")
                for i, col := range result.Columns {
                    fmt.Printf("%s", col)
                    if i < len(result.Columns)-1 {
                        fmt.Print(" | ")
                    }
                }
                fmt.Print(" |\n")

                // Print separator
                fmt.Print("|")
                for _, col := range result.Columns {
                    fmt.Print(strings.Repeat("-", len(col)+2))
                    fmt.Print("|")
                }
                fmt.Print("\n")

                // Print records
                for _, record := range result.Records {
                    fmt.Print("| ")
                    for i, col := range result.Columns {
                        value := record.Columns[col]
                        if value == nil {
                            fmt.Print("NULL")
                        } else {
                            fmt.Printf("%v", value)
                        }
                        if i < len(result.Columns)-1 {
                            fmt.Print(" | ")
                        }
                    }
                    fmt.Print(" |\n")
                }
            } else if result.Message != "" {
                fmt.Println(result.Message)
            }
            
            fmt.Println()
        }
        
        return
    }

    // Create a scanner to read input
    scanner := bufio.NewScanner(os.Stdin)
    var currentCommand string

    fmt.Print("> ")
    for scanner.Scan() {
        line := scanner.Text()

        // Skip empty lines
        if strings.TrimSpace(line) == "" {
            fmt.Print("> ")
            continue
        }

        // Append line to current command
        if currentCommand != "" {
            currentCommand += "\n"
        }
        currentCommand += line

        // Check if command is complete (ends with semicolon)
        if !strings.HasSuffix(strings.TrimSpace(currentCommand), ";") {
            fmt.Print("... ")
            continue
        }

        // Parse and execute command
        stmt, err := sql.Parse(currentCommand)
        if err != nil {
            fmt.Printf("Error parsing command '%s': %v\n", currentCommand, err)
            currentCommand = ""
            fmt.Print("> ")
            continue
        }

        // Skip empty statements (comments)
        if stmt == nil {
            currentCommand = ""
            fmt.Print("> ")
            continue
        }

        // Execute statement
        result, err := database.Execute(stmt)
        if err != nil {
            fmt.Printf("Error executing command '%s': %v\n", currentCommand, err)
            currentCommand = ""
            fmt.Print("> ")
            continue
        }

        // Print result
        if result.IsSelect {
            // Print table header
            fmt.Print("| ")
            for i, col := range result.Columns {
                fmt.Printf("%s", col)
                if i < len(result.Columns)-1 {
                    fmt.Print(" | ")
                }
            }
            fmt.Print(" |\n")

            // Print separator
            fmt.Print("|")
            for _, col := range result.Columns {
                fmt.Print(strings.Repeat("-", len(col)+2))
                fmt.Print("|")
            }
            fmt.Print("\n")

            // Print records
            for _, record := range result.Records {
                fmt.Print("| ")
                for i, col := range result.Columns {
                    value := record.Columns[col]
                    if value == nil {
                        fmt.Print("NULL")
                    } else {
                        fmt.Printf("%v", value)
                    }
                    if i < len(result.Columns)-1 {
                        fmt.Print(" | ")
                    }
                }
                fmt.Print(" |\n")
            }
        } else if result.Message != "" {
            fmt.Println(result.Message)
        }

        currentCommand = ""
        fmt.Print("> ")
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input: %v\n", err)
    }

    fmt.Println("\nINFO: \nExiting due to EOF. Goodbye!")
}

func printWelcome() {
    welcome := `
······································································
: ________  ________  ___       ___  ________  ___  ___  _________   :
:|\   ____\|\   __  \|\  \     |\  \|\   ____\|\  \|\  \|\___   ___\ :
:\ \  \___|\ \  \|\  \ \  \    \ \  \ \  \___|\ \  \\\  \|___ \  \_| :
: \ \_____  \ \  \\\  \ \  \    \ \  \ \  \  __\ \   __  \   \ \  \  :
:  \|____|\  \ \  \\\  \ \  \____\ \  \ \  \|\  \ \  \ \  \   \ \  \ :
:    ____\_\  \ \_____  \ \_______\ \__\ \_______\ \__\ \__\   \ \__\:
:   |\_________\|___| \__\|_______|\|__|\|_______|\|__|\|__|    \|__|:
:   \|_________|     \|__|                                           :
······································································
`
    fmt.Println(welcome)
    fmt.Println("Welcome to SQLight! Type 'help' for usage information.")
    fmt.Println("Using database file: database.json\n")
}

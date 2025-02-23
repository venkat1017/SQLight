package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sqlite-clone/pkg/db"
	"sqlite-clone/pkg/logger"
	"sqlite-clone/pkg/sql"
	"strings"
	"syscall"
)

const banner = `
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
const helpText = `
Available Commands:
  CREATE TABLE users (id INTEGER, name TEXT, email TEXT)  - Create a new table
  INSERT INTO users VALUES (1, 'Alice', 'alice@email.com') - Insert a record
  SELECT * FROM users                                      - Select all records
  SELECT * FROM users WHERE id = 1                         - Select by ID
  SELECT * FROM users WHERE name = 'Alice'                 - Select by name
  UPDATE users SET name = 'Alice Smith' WHERE id = 1       - Update records
  DELETE FROM users WHERE id = 1                           - Delete records
  exit                                                     - Exit the program
  help                                                     - Show this help message

Tips:
  - String values must be enclosed in single quotes ('value')
  - Commands are case-insensitive
  - Commands can end with or without a semicolon (;)
  - Use Ctrl+C to exit safely at any time
`

func main() {
	// Parse command line flags
	debug := flag.Bool("debug", false, "Enable debug logging")
	dbFile := flag.String("db", "database.json", "Database file path")
	flag.Parse()

	// Set up logging
	logger.SetDebug(*debug)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize database
	database := db.NewDatabase(*dbFile)
	defer database.Save()

	reader := bufio.NewReader(os.Stdin)

	// Print welcome message
	fmt.Println(banner)
	fmt.Println("Welcome to SQLite Clone! Type 'help' for usage information.")
	fmt.Println("Using database file:", *dbFile)

	// Start a goroutine to handle signals
	go func() {
		<-sigChan
		logger.Infof("\nReceived interrupt signal. Saving and exiting...")
		database.Save()
		os.Exit(0)
	}()

	// Main input loop
	for {
		fmt.Print("\n> ") // Cleaner prompt
		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				logger.Infof("\nExiting due to EOF. Goodbye!")
				break
			}
			logger.Errorf("Error reading input: %v", err)
			continue
		}

		// Trim the input
		input = strings.TrimSpace(input)

		// Skip empty input
		if input == "" {
			continue
		}

		// Handle special commands
		switch strings.ToLower(input) {
		case "exit", "quit":
			fmt.Println("Saving and exiting. Goodbye!")
			return
		case "help":
			fmt.Println(helpText)
			continue
		}

		// Parse and execute the SQL command
		stmt, err := sql.ParseSQL(input)
		if err != nil {
			logger.Errorf("Error parsing command '%s': %v", input, err)
			continue
		}

		if err := stmt.Exec(database); err != nil {
			logger.Errorf("Error executing command '%s': %v", input, err)
		}

		logger.Debugf("Command executed successfully: %s", input)
	}
}

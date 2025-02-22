package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	// Debug controls whether debug messages are printed
	Debug bool

	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

// SetDebug enables or disables debug logging
func SetDebug(enabled bool) {
	Debug = enabled
}

// Debugf prints debug messages if debug mode is enabled
func Debugf(format string, args ...interface{}) {
	if Debug {
		logger.Printf("DEBUG: "+format, args...)
	}
}

// Infof prints info messages
func Infof(format string, args ...interface{}) {
	logger.Printf("INFO: "+format, args...)
}

// Errorf prints error messages
func Errorf(format string, args ...interface{}) {
	logger.Printf("ERROR: "+format, args...)
}

// PrintTable prints a table in a formatted way
func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		fmt.Println("No data to display")
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print top border
	printBorder(widths)

	// Print headers
	fmt.Print("|")
	for i, h := range headers {
		fmt.Printf(" %-*s |", widths[i], h)
	}
	fmt.Println()

	// Print separator
	printBorder(widths)

	// Print rows
	for _, row := range rows {
		fmt.Print("|")
		for i, cell := range row {
			fmt.Printf(" %-*s |", widths[i], cell)
		}
		fmt.Println()
	}

	// Print bottom border
	printBorder(widths)
}

func printBorder(widths []int) {
	fmt.Print("+")
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2) + "+")
	}
	fmt.Println()
}

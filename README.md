# SQLite Clone in Go

A lightweight SQLite clone implemented in Go that supports basic SQL operations. This project demonstrates how to create a simple but functional database system with support for CRUD operations (Create, Read, Update, Delete) and data type validation.

## Features

- **SQL Command Support**:
  - `CREATE TABLE` - Create tables with specified columns and data types
  - `INSERT INTO` - Insert records into tables
  - `SELECT` - Query records with WHERE clause support
  - `UPDATE` - Update existing records
  - `DELETE` - Delete records from tables

- **Data Types**:
  - `INTEGER` - Whole numbers
  - `TEXT` - String values
  - `BOOLEAN` - True/False values
  - `DATETIME` - Date and time values

- **Additional Features**:
  - Case-insensitive SQL commands
  - Persistent storage using JSON
  - Data type validation
  - B-tree index for efficient lookups
  - Debug logging support
  - Clean command-line interface

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/sqlight.git
cd sqlight
```

2. Build the project:
```bash
go build -o sqlight ./cmd/main.go
```

## Usage

Run the program:
```bash
# Basic usage
./sqlight

# With debug logging enabled
./sqlight -debug

# With custom database file
./sqlight -db mydb.json
```

### Example SQL Commands

1. Create a table:
```sql
CREATE TABLE users (id INTEGER, name TEXT, email TEXT);
```

2. Insert records:
```sql
INSERT INTO users VALUES (1, 'Alice', 'alice@email.com');
INSERT INTO users VALUES (2, 'Bob', 'bob@email.com');
```

3. Query records:
```sql
-- Select all records
SELECT * FROM users;

-- Select with WHERE clause
SELECT * FROM users WHERE id = 1;
SELECT * FROM users WHERE name = 'Alice';
```

4. Update records:
```sql
UPDATE users SET name = 'Alice Smith' WHERE id = 1;
```

5. Delete records:
```sql
DELETE FROM users WHERE id = 2;
```

### Additional Commands
- Type `help` to see available commands and usage tips
- Type `exit` or `quit` to safely exit the program
- Use Ctrl+C to exit safely at any time

## Project Structure

```
sqlight/
├── cmd/
│   └── main.go           # Main application entry point
├── pkg/
│   ├── db/
│   │   ├── database.go   # Database operations
│   │   ├── table.go      # Table operations
│   │   └── btree.go      # B-tree implementation
│   ├── sql/
│   │   └── parser.go     # SQL parser
│   ├── types/
│   │   └── datatypes/    # Data type implementations
│   └── logger/
│       └── logger.go     # Logging utilities
└── README.md
```

## Development

### Prerequisites
- Go 1.16 or later
- Any text editor or IDE

### Running Tests
```bash
go test ./...
```

### Debug Mode
Run with the `-debug` flag to enable detailed logging:
```bash
./sqlight -debug
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by [SQLite](https://sqlite.org/)
- Built using Go's standard library
- Uses B-tree data structure for efficient indexing

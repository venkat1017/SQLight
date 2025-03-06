# SQLight - A Modern SQLite Clone with Web Interface

A lightweight SQLite clone implemented in Go that supports basic SQL operations through both a command-line interface and a modern web interface. This project demonstrates creating a simple but functional database system with persistent storage and a clean UI.

## Features

- **Dual Interface Support**:
  - Modern Web Interface with real-time query execution
  - Traditional Command Line Interface
  
- **SQL Command Support**:
  - `CREATE TABLE` - Create tables with specified columns and data types
  - `INSERT INTO` - Insert records into tables
  - `SELECT` - Query records from tables
  - More commands coming soon!

- **Data Types**:
  - `INTEGER` - Whole numbers
  - `TEXT` - String values
  - More types coming soon!

- **Additional Features**:
  - Case-insensitive SQL commands
  - Persistent storage using JSON
  - Data type validation
  - Clean, modern web interface
  - Real-time query execution
  - Table list sidebar
  - Success/Error messages
  - Keyboard shortcuts (Ctrl+Enter/Cmd+Enter to run queries)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/sqlight.git
cd sqlight
```

2. Build the project:
```bash
# Build CLI version
go build -o sqlight ./cmd/main.go

# Build web version
go build -o sqlightweb ./web/main.go
```

## Usage

### Web Interface

1. Start the web server:
```bash
./sqlightweb
# Or run directly with Go
go run web/main.go
```

2. Open your browser and visit:
```
http://localhost:8081
```

3. Use the web interface to:
- Write and execute SQL queries
- View table list in the sidebar
- Click on tables to auto-fill SELECT queries
- See success/error messages
- View query results in a formatted table

### Command Line Interface

Run the CLI version:
```bash
# Basic usage
./sqlight

# With custom database file
./sqlight -db mydb.json
```

### Example SQL Commands

1. Create a table:
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE
);
```

2. Insert records:
```sql
INSERT INTO users (id, name, email) VALUES (1, 'John Doe', 'john@example.com');
```

3. Query records:
```sql
-- Select all records
SELECT * FROM users;
```

## Project Structure

```
sqlight/
├── cmd/
│   └── main.go           # CLI application entry
├── web/
│   ├── main.go          # Web server entry
│   └── static/          # Web interface files
│       ├── index.html   # Main HTML page
│       ├── styles.css   # CSS styles
│       └── script.js    # Frontend JavaScript
├── pkg/
│   ├── db/             # Database implementation
│   │   ├── database.go # Database operations
│   │   ├── table.go    # Table operations
│   │   └── cursor.go   # Record cursor
│   ├── sql/            # SQL parsing
│   │   └── parser.go   # SQL parser
│   └── interfaces/     # Core interfaces
│       └── interfaces.go
└── README.md
```

## Development

### Prerequisites
- Go 1.16 or later
- Modern web browser for web interface
- Any text editor or IDE

### Running Tests
```bash
go test ./...
```

### Browser Support
The web interface works best with:
- Chrome/Edge (latest versions)
- Firefox (latest version)
- Safari (latest version)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Areas for Contribution
- Additional SQL command support (UPDATE, DELETE)
- More data types (FLOAT, DATETIME, etc.)
- Improved SQL parsing
- Query optimization
- Additional indexes
- UI/UX improvements
- Documentation
- Tests

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by [SQLite](https://sqlite.org/)
- Built using Go's standard library
- Modern web interface using vanilla JavaScript

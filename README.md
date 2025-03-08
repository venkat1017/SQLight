# SQLight

<div align="center">

![SQLight Logo](https://img.shields.io/badge/SQLight-A%20Modern%20SQLite%20Clone-blue?style=for-the-badge)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16%2B-00ADD8.svg)](https://golang.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Last Commit](https://img.shields.io/github/last-commit/venkat1017/litesqlite)](https://github.com/venkat1017/litesqlite/commits/main)

**A lightweight SQLite clone with a modern web interface and CLI**

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Documentation](#documentation) • [Contributing](#contributing) • [License](#license)

</div>

## 📋 Overview

SQLight is a lightweight SQLite clone implemented in Go that provides a simple yet functional database system with persistent storage. It features both a modern web interface and a traditional command-line interface, making it versatile for different use cases.

This project demonstrates core database concepts including SQL parsing, query execution, transaction management, and persistent storage while maintaining a clean, user-friendly interface.

<div align="center">
  <img src="https://img.shields.io/badge/SQLight-Web%20Interface%20Demo-blue?style=for-the-badge" alt="Web Interface Demo">
</div>

## ✨ Features

### 🖥️ Dual Interface Support
- **Modern Web Interface** with real-time query execution and interactive table browsing
- **Traditional Command Line Interface** for script-based and terminal operations

### 📊 SQL Command Support
- `CREATE TABLE` - Create tables with specified columns and data types
- `INSERT INTO` - Insert records into tables
- `SELECT` - Query records with support for WHERE clauses and column selection
- `DELETE` - Remove records with WHERE clause filtering
- More commands coming soon!

### 🔄 Data Types
- `INTEGER` - Whole numbers
- `TEXT` - String values
- More types coming soon!

### 🛠️ Advanced Features
- **Transaction Support** for atomic operations
- **Case-insensitive** SQL command and table/column name handling
- **WHERE Clause Support** with multiple conditions using AND
- **String Value Handling** with support for both single and double quotes
- **Persistent Storage** using JSON
- **Data Type Validation** for integrity
- **Error Handling** for non-existent tables/columns
- **B-tree Implementation** for efficient data storage and retrieval

### 🎨 Web Interface Features
- Clean, modern UI with dark/light mode support
- Real-time query execution
- Interactive table list sidebar
- Success/Error messages with detailed feedback
- Keyboard shortcuts (Ctrl+Enter/Cmd+Enter to run queries)
- Responsive design for desktop and tablet use

## 🚀 Installation

### Prerequisites
- Go 1.16 or later
- Modern web browser for web interface
- Git (for cloning the repository)

### Quick Start

1. **Clone the repository**:
```bash
git clone https://github.com/venkat1017/litesqlite.git
cd litesqlite
```

2. **Build the project**:
```bash
# Build CLI version
go build -o sqlight ./cmd/main.go

# Build web version
go build -o sqlightweb ./web/main.go
```

3. **Run directly with Go** (alternative to building):
```bash
# Run CLI version
go run cmd/main.go

# Run web version
go run web/main.go
```

## 🖱️ Usage

### Web Interface

1. **Start the web server**:
```bash
./sqlightweb
# Or run directly with Go
go run web/main.go
```

2. **Open your browser** and visit:
```
http://localhost:8081
```

3. **Use the web interface to**:
- Write and execute SQL queries
- View table list in the sidebar
- Click on tables to auto-fill SELECT queries
- See detailed success/error messages
- View query results in a formatted table

### Command Line Interface

**Run the CLI version**:
```bash
# Basic usage
./sqlight

# With custom database file
./sqlight -db mydb.json
```

## 📝 Example SQL Commands

### Create a Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE
);
```

### Insert Records
```sql
INSERT INTO users (id, name, email) VALUES (1, 'John Doe', 'john@example.com');
INSERT INTO users (id, name, email) VALUES (2, 'Jane Smith', 'jane@example.com');
```

### Query Records
```sql
-- Select all records
SELECT * FROM users;

-- Select specific columns
SELECT id, name FROM users;

-- Select with WHERE clause
SELECT * FROM users WHERE id = 1;

-- Select with multiple conditions
SELECT * FROM users WHERE id > 0 AND name = 'John Doe';
```

### Delete Records
```sql
-- Delete specific records
DELETE FROM users WHERE id = 1;

-- Delete with multiple conditions
DELETE FROM users WHERE id > 5 AND name = 'Test User';

-- Delete all records from a table
DELETE FROM users;
```

## 📁 Project Structure

```
sqlight/
├── cmd/                  # Command-line application
│   └── main.go           # CLI entry point
├── web/                  # Web server application
│   ├── main.go           # Web server entry point
│   └── static/           # Web interface files
│       ├── index.html    # Main HTML page
│       ├── styles.css    # CSS styles
│       └── script.js     # Frontend JavaScript
├── pkg/                  # Core packages
│   ├── db/               # Database implementation
│   │   ├── database.go   # Database operations
│   │   ├── table.go      # Table operations
│   │   ├── btree.go      # B-tree implementation
│   │   └── cursor.go     # Record cursor
│   ├── sql/              # SQL parsing
│   │   └── parser.go     # SQL parser
│   └── interfaces/       # Core interfaces
│       └── interfaces.go # Interface definitions
├── examples/             # Example usage and demos
├── tests/                # Test suite
├── go.mod                # Go module definition
└── README.md             # Project documentation
```

## 📚 Documentation

### Architecture

SQLight follows a layered architecture:

1. **Interface Layer** - Web UI and CLI for user interaction
2. **SQL Parser** - Converts SQL strings into structured statements
3. **Query Executor** - Processes statements and performs operations
4. **Storage Engine** - Manages data persistence and retrieval
5. **B-tree Implementation** - Provides efficient data storage and access

### Performance Considerations

- SQLight uses a B-tree implementation for efficient data access
- JSON-based persistence provides a balance of simplicity and performance
- In-memory operations for speed with periodic persistence for durability

## 🧪 Development

### Running Tests
```bash
go test ./...
```

### Debugging
```bash
# Run with verbose logging
go run cmd/main.go -v
```

### Browser Support
The web interface works best with:
- Chrome/Edge (latest versions)
- Firefox (latest version)
- Safari (latest version)

## 👥 Contributing

We welcome contributions from the community! Here's how you can help:

1. **Fork** the repository
2. **Create** your feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add some amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Areas for Contribution
- Additional SQL command support (UPDATE, JOIN operations)
- More data types (FLOAT, DATETIME, BOOLEAN, etc.)
- Improved SQL parsing and validation
- Query optimization and execution planning
- Additional indexing strategies
- UI/UX improvements
- Documentation enhancements
- Test coverage expansion

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [SQLite](https://sqlite.org/)
- Built using Go's standard library
- Modern web interface using vanilla JavaScript
- Thanks to all contributors who have helped shape this project

---

<div align="center">
  
**[⬆ Back to Top](#sqlight)**

Made with ❤️ by the SQLight team

</div>

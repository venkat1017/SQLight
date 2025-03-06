-- Basic SQL queries example
-- Run these queries to get started with SQLight

-- Create a users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE
);

-- Insert some sample data
INSERT INTO users (id, name, email) VALUES (1, 'John Doe', 'john@example.com');
INSERT INTO users (id, name, email) VALUES (2, 'Jane Smith', 'jane@example.com');
INSERT INTO users (id, name, email) VALUES (3, 'Bob Wilson', 'bob@example.com');

-- Query all users
SELECT * FROM users;

-- Query specific user
SELECT * FROM users WHERE id = 1;

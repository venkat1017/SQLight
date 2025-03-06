-- Blog database schema example
-- This example shows how to create a simple blog database structure

-- Create authors table
CREATE TABLE authors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    bio TEXT
);

-- Create posts table
CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TEXT NOT NULL
);

-- Create comments table
CREATE TABLE comments (
    id INTEGER PRIMARY KEY,
    post_id INTEGER NOT NULL,
    author_name TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Insert sample data
INSERT INTO authors (id, name, email, bio) 
VALUES (1, 'John Doe', 'john@blog.com', 'Tech writer and software developer');

INSERT INTO posts (id, title, content, author_id, created_at)
VALUES (1, 'Getting Started with SQLight', 'SQLight is a lightweight database...', 1, '2023-11-15');

INSERT INTO comments (id, post_id, author_name, content, created_at)
VALUES (1, 1, 'Jane Smith', 'Great article!', '2023-11-15');

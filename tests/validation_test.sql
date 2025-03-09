-- SQLight Validation Test Script
-- This script tests all supported SQL commands including DELETE

-- Create a test table
CREATE TABLE test_users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT,
    age INTEGER
);

-- Insert test data
INSERT INTO test_users (id, name, email, age) VALUES (1, 'John Doe', 'john@example.com', 30);
INSERT INTO test_users (id, name, email, age) VALUES (2, 'Jane Smith', 'jane@example.com', 25);
INSERT INTO test_users (id, name, email, age) VALUES (3, 'Bob Johnson', 'bob@example.com', 40);
INSERT INTO test_users (id, name, email, age) VALUES (4, 'Alice Brown', 'alice@example.com', 35);
INSERT INTO test_users (id, name, email, age) VALUES (5, 'Test User', 'test@example.com', 20);

-- Select all records to verify insertion
SELECT * FROM test_users;

-- Test SELECT with WHERE clause
SELECT * FROM test_users WHERE age > 30;

-- Test SELECT with multiple conditions
SELECT * FROM test_users WHERE age > 20 AND name = 'Jane Smith';

-- Test DELETE with single condition
DELETE FROM test_users WHERE id = 5;

-- Verify the deletion worked
SELECT * FROM test_users;

-- Test DELETE with multiple conditions using AND
DELETE FROM test_users WHERE age > 30 AND name = 'Bob Johnson';

-- Verify the deletion with multiple conditions worked
SELECT * FROM test_users;

-- Create another table to test case-insensitivity
CREATE TABLE Test_Items (
    ItemId INTEGER PRIMARY KEY,
    ItemName TEXT,
    Price INTEGER
);

-- Insert data with mixed case
INSERT INTO test_items (ItemId, ItemName, Price) VALUES (1, 'Laptop', 1000);
INSERT INTO test_items (ItemId, ItemName, Price) VALUES (2, 'Phone', 500);

-- Test case-insensitive SELECT
SELECT * FROM TEST_ITEMS;

-- Test case-insensitive DELETE
DELETE FROM TEST_items WHERE itemid = 2;

-- Verify case-insensitive DELETE worked
SELECT * FROM test_items;

-- Test string values in WHERE conditions with different quote styles
INSERT INTO test_users (id, name, email, age) VALUES (6, 'Quote Test', 'quotes@example.com', 45);
SELECT * FROM test_users WHERE name = 'Quote Test';
SELECT * FROM test_users WHERE name = "Quote Test";

-- Test DELETE with string value in WHERE condition
DELETE FROM test_users WHERE email = 'quotes@example.com';

-- Verify string value DELETE worked
SELECT * FROM test_users;

-- SQLite Clone Demo
-- This file demonstrates all the functionality of the SQLite clone

-- 1. CREATE TABLE with various constraints
CREATE TABLE employees (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    department TEXT NOT NULL,
    salary FLOAT
);

-- Display table structure
DESCRIBE employees;

-- 2. INSERT with various data types
-- Basic inserts
INSERT INTO employees (id, name, department, salary) VALUES (1, 'John Doe', 'Engineering', 85000.50);
INSERT INTO employees (id, name, department, salary) VALUES (2, 'Jane Smith', 'Marketing', 75000);
INSERT INTO employees (id, name, department, salary) VALUES (3, 'Bob Johnson', 'Finance', 90000);
INSERT INTO employees (id, name, department, salary) VALUES (4, 'Alice Brown', 'Engineering', 82000);
INSERT INTO employees (id, name, department, salary) VALUES (5, 'Charlie Davis', 'HR', 65000);

-- Show all records
SELECT * FROM employees;

-- 3. Constraint testing
-- PRIMARY KEY constraint (should fail)
INSERT INTO employees (id, name, department, salary) VALUES (1, 'Duplicate ID', 'Research', 95000);

-- UNIQUE constraint (should fail)
INSERT INTO employees (id, name, department, salary) VALUES (6, 'John Doe', 'Research', 95000);

-- NOT NULL constraint (should fail)
INSERT INTO employees (id, name, salary) VALUES (7, 'Missing Department', 70000);

-- 4. SELECT with various conditions
-- Select all columns with WHERE clause
SELECT * FROM employees WHERE department = 'Engineering';

-- Select specific columns
SELECT name, salary FROM employees;

-- Select with numeric comparison
SELECT * FROM employees WHERE salary > 80000;

-- Select with string condition
SELECT * FROM employees WHERE name = 'Jane Smith';

-- Case insensitive column and table names
SELECT NAME, DEPARTMENT FROM EMPLOYEES WHERE DEPARTMENT = 'Marketing';

-- 5. DELETE operations
-- Delete with WHERE clause
DELETE FROM employees WHERE id = 5;
SELECT * FROM employees;

-- Delete with string condition
DELETE FROM employees WHERE name = 'Bob Johnson';
SELECT * FROM employees;

-- Delete with multiple conditions
INSERT INTO employees (id, name, department, salary) VALUES (6, 'Test User', 'Test Dept', 50000);
INSERT INTO employees (id, name, department, salary) VALUES (7, 'Test User2', 'Test Dept', 55000);
SELECT * FROM employees;
DELETE FROM employees WHERE department = 'Test Dept' AND salary = 50000;
SELECT * FROM employees;

-- 6. Transaction support
-- Begin a transaction
BEGIN TRANSACTION;

-- Make some changes
INSERT INTO employees (id, name, department, salary) VALUES (8, 'Transaction Test', 'Legal', 72000);
DELETE FROM employees WHERE id = 7;
SELECT * FROM employees;

-- Rollback the transaction
ROLLBACK;

-- Verify changes were rolled back
SELECT * FROM employees;

-- Begin another transaction
BEGIN TRANSACTION;

-- Make some changes
INSERT INTO employees (id, name, department, salary) VALUES (8, 'Committed User', 'Legal', 72000);
DELETE FROM employees WHERE id = 7;
SELECT * FROM employees;

-- Commit the transaction
COMMIT;

-- Verify changes were committed
SELECT * FROM employees;

-- 7. DROP TABLE
-- Create a temporary table
CREATE TABLE temp_table (id INT, name TEXT);
INSERT INTO temp_table (id, name) VALUES (1, 'Temporary');
SELECT * FROM temp_table;

-- Drop the table
DROP TABLE temp_table;

-- 8. DELETE without WHERE clause
CREATE TABLE test_delete_all (id INT, name TEXT);
INSERT INTO test_delete_all (id, name) VALUES (1, 'Delete Me');
INSERT INTO test_delete_all (id, name) VALUES (2, 'Delete Me Too');
SELECT * FROM test_delete_all;

-- Delete all records
DELETE FROM test_delete_all;
SELECT * FROM test_delete_all;

-- Final display of main table
SELECT * FROM employees;

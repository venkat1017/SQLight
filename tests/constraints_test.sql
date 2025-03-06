-- Constraints Test SQL File
-- This file tests various SQL constraints including PRIMARY KEY, UNIQUE, NOT NULL, and FOREIGN KEY

-- Clean up any existing tables
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;

-- Create departments table with PRIMARY KEY constraint
CREATE TABLE departments (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    location TEXT
);

-- Create employees table with PRIMARY KEY, UNIQUE, NOT NULL, and FOREIGN KEY constraints
CREATE TABLE employees (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    salary REAL,
    department_id INTEGER,
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- Describe tables to verify constraints
DESCRIBE departments;
DESCRIBE employees;

-- Test 1: Insert valid data into departments
INSERT INTO departments (id, name, location) VALUES (1, 'Engineering', 'Building A');
INSERT INTO departments (id, name, location) VALUES (2, 'Marketing', 'Building B');
INSERT INTO departments (id, name, location) VALUES (3, 'Finance', 'Building C');

-- Test 2: Insert valid data into employees
INSERT INTO employees (id, email, name, salary, department_id) VALUES (1, 'john@example.com', 'John Doe', 75000, 1);
INSERT INTO employees (id, email, name, salary, department_id) VALUES (2, 'jane@example.com', 'Jane Smith', 85000, 1);
INSERT INTO employees (id, email, name, salary, department_id) VALUES (3, 'bob@example.com', 'Bob Johnson', 65000, 2);

-- Test 3: PRIMARY KEY constraint - Should fail (duplicate id)
INSERT INTO departments (id, name, location) VALUES (1, 'Research', 'Building D');

-- Test 4: UNIQUE constraint - Should fail (duplicate name)
INSERT INTO departments (id, name, location) VALUES (4, 'Engineering', 'Building E');

-- Test 5: NOT NULL constraint - Should fail (null name)
INSERT INTO departments (id, location) VALUES (5, 'Building F');

-- Test 6: FOREIGN KEY constraint - Should fail (invalid department_id)
INSERT INTO employees (id, email, name, salary, department_id) VALUES (4, 'alice@example.com', 'Alice Brown', 70000, 99);

-- Test 7: UNIQUE constraint - Should fail (duplicate email)
INSERT INTO employees (id, email, name, salary, department_id) VALUES (5, 'john@example.com', 'John Smith', 60000, 3);

-- Test 8: NOT NULL constraint - Should fail (null email)
INSERT INTO employees (id, name, salary, department_id) VALUES (6, 'Sarah Wilson', 90000, 3);

-- Query data to verify successful insertions
SELECT * FROM departments;
SELECT * FROM employees;

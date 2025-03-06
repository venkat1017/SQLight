-- Debug constraints
CREATE TABLE test_pk (id INTEGER PRIMARY KEY, name TEXT);
DESCRIBE TABLE test_pk;

-- Test with explicit SQL commands to check if constraints are being set
CREATE TABLE test_constraints (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    dept_id INTEGER
);
DESCRIBE TABLE test_constraints;

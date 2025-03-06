-- Simple constraints test
CREATE TABLE test_dept (id INTEGER PRIMARY KEY, name TEXT UNIQUE NOT NULL);
DESCRIBE test_dept;

-- Insert valid data
INSERT INTO test_dept (id, name) VALUES (1, 'Engineering');
INSERT INTO test_dept (id, name) VALUES (2, 'Marketing');
SELECT * FROM test_dept;

-- Test PRIMARY KEY constraint (should fail)
INSERT INTO test_dept (id, name) VALUES (1, 'Research');
SELECT * FROM test_dept;

-- Test UNIQUE constraint (should fail)
INSERT INTO test_dept (id, name) VALUES (3, 'Engineering');
SELECT * FROM test_dept;

-- Test NOT NULL constraint (should fail)
INSERT INTO test_dept (id) VALUES (4);
SELECT * FROM test_dept;

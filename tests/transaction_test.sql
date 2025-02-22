-- Create a test table
CREATE TABLE users (id INTEGER, name TEXT);

-- Start a transaction
BEGIN TRANSACTION;

-- Insert some records
INSERT INTO users (id, name) VALUES (1, 'Alice');
INSERT INTO users (id, name) VALUES (2, 'Bob');

-- View the records (should show both Alice and Bob)
SELECT * FROM users;

-- Commit the transaction
COMMIT;

-- Start another transaction
BEGIN TRANSACTION;

-- Insert more records
INSERT INTO users (id, name) VALUES (3, 'Charlie');
INSERT INTO users (id, name) VALUES (4, 'David');

-- View all records (should show all 4 users)
SELECT * FROM users;

-- Rollback this transaction
ROLLBACK;

-- View records again (should only show Alice and Bob)
SELECT * FROM users;

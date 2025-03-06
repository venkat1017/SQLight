-- Inventory management system example
-- Demonstrates a simple inventory tracking system

-- Create products table
CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL
);

-- Create categories table
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create product_categories table for many-to-many relationship
CREATE TABLE product_categories (
    product_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL
);

-- Insert sample categories
INSERT INTO categories (id, name) VALUES (1, 'Electronics');
INSERT INTO categories (id, name) VALUES (2, 'Books');
INSERT INTO categories (id, name) VALUES (3, 'Clothing');

-- Insert sample products
INSERT INTO products (id, name, description, price, quantity)
VALUES (1, 'Laptop', 'High-performance laptop', 999, 10);

INSERT INTO products (id, name, description, price, quantity)
VALUES (2, 'T-Shirt', 'Cotton t-shirt', 20, 100);

-- Associate products with categories
INSERT INTO product_categories (product_id, category_id) VALUES (1, 1);
INSERT INTO product_categories (product_id, category_id) VALUES (2, 3);

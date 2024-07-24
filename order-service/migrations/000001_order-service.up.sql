CREATE TABLE orders (
    product_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    total_amount FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
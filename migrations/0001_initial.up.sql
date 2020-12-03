CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    customer_ref_id VARCHAR(100),
    product_name VARCHAR(128),
    current_price INT,
    original_price INT,
    source VARCHAR(100),
    is_deleted BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

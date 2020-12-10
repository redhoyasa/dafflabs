CREATE TABLE wishes (
    id VARCHAR(100) PRIMARY KEY,
    customer_ref_id VARCHAR(100) NOT NULL,
    product_name VARCHAR(128) NOT NULL,
    current_price INT NOT NULL,
    original_price INT NOT NULL,
    discount_rate INT,
    source VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    UNIQUE (customer_ref_id, source)
);

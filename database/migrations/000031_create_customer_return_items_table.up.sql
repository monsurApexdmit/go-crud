CREATE TABLE IF NOT EXISTS customer_return_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    return_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NULL,
    product_name VARCHAR(255) NOT NULL,
    variant_id BIGINT UNSIGNED NULL,
    variant_name VARCHAR(255) NULL,
    quantity INT NOT NULL DEFAULT 1,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (return_id) REFERENCES customer_returns (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE SET NULL,
    FOREIGN KEY (variant_id) REFERENCES product_variants (id) ON DELETE SET NULL,
    INDEX idx_customer_return_items_return_id (return_id),
    INDEX idx_customer_return_items_product_id (product_id)
);

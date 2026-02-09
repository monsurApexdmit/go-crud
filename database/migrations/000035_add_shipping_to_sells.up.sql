ALTER TABLE sells
ADD COLUMN shipping_address_id BIGINT UNSIGNED NULL AFTER customer_name,
ADD COLUMN shipping_method VARCHAR(50) NULL AFTER shipping_cost,
ADD COLUMN payment_status ENUM('pending', 'paid', 'partially_paid', 'refunded', 'failed') DEFAULT 'pending' AFTER status,
ADD COLUMN fulfillment_status ENUM('unfulfilled', 'processing', 'shipped', 'delivered', 'cancelled') DEFAULT 'unfulfilled' AFTER payment_status,
ADD COLUMN tracking_number VARCHAR(100) NULL AFTER fulfillment_status,
ADD COLUMN carrier VARCHAR(100) NULL AFTER tracking_number,
ADD COLUMN shipped_at TIMESTAMP NULL AFTER carrier,
ADD COLUMN delivered_at TIMESTAMP NULL AFTER shipped_at,
ADD FOREIGN KEY (shipping_address_id) REFERENCES shipping_addresses (id) ON DELETE SET NULL,
ADD INDEX idx_sells_shipping_address_id (shipping_address_id),
ADD INDEX idx_sells_payment_status (payment_status),
ADD INDEX idx_sells_fulfillment_status (fulfillment_status),
ADD INDEX idx_sells_tracking_number (tracking_number);

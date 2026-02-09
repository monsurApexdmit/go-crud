-- Add shipping address fields directly to sells table
-- This creates a snapshot of the shipping address at order time
-- Industry standard: Preserve historical data even if customer updates/deletes their address

ALTER TABLE sells
ADD COLUMN shipping_full_name VARCHAR(255) NULL AFTER shipping_address_id,
ADD COLUMN shipping_phone VARCHAR(20) NULL AFTER shipping_full_name,
ADD COLUMN shipping_email VARCHAR(255) NULL AFTER shipping_phone,
ADD COLUMN shipping_address_line1 VARCHAR(255) NULL AFTER shipping_email,
ADD COLUMN shipping_address_line2 VARCHAR(255) NULL AFTER shipping_address_line1,
ADD COLUMN shipping_city VARCHAR(100) NULL AFTER shipping_address_line2,
ADD COLUMN shipping_state VARCHAR(100) NULL AFTER shipping_city,
ADD COLUMN shipping_postal_code VARCHAR(20) NULL AFTER shipping_state,
ADD COLUMN shipping_country VARCHAR(100) NULL AFTER shipping_postal_code,
ADD COLUMN shipping_address_type ENUM('home', 'office', 'other') NULL AFTER shipping_country;

-- Add indexes for common queries
CREATE INDEX idx_sells_shipping_city ON sells(shipping_city);
CREATE INDEX idx_sells_shipping_postal_code ON sells(shipping_postal_code);
CREATE INDEX idx_sells_shipping_country ON sells(shipping_country);

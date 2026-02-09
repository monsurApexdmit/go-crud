-- Rollback: Remove shipping address snapshot fields from sells table

DROP INDEX idx_sells_shipping_city ON sells;
DROP INDEX idx_sells_shipping_postal_code ON sells;
DROP INDEX idx_sells_shipping_country ON sells;

ALTER TABLE sells
DROP COLUMN shipping_full_name,
DROP COLUMN shipping_phone,
DROP COLUMN shipping_email,
DROP COLUMN shipping_address_line1,
DROP COLUMN shipping_address_line2,
DROP COLUMN shipping_city,
DROP COLUMN shipping_state,
DROP COLUMN shipping_postal_code,
DROP COLUMN shipping_country,
DROP COLUMN shipping_address_type;

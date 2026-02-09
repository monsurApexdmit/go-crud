ALTER TABLE sells
DROP FOREIGN KEY sells_ibfk_2,
DROP COLUMN delivered_at,
DROP COLUMN shipped_at,
DROP COLUMN carrier,
DROP COLUMN tracking_number,
DROP COLUMN fulfillment_status,
DROP COLUMN payment_status,
DROP COLUMN shipping_method,
DROP COLUMN shipping_address_id;

ALTER TABLE product_variants
  DROP INDEX sku,
  DROP INDEX barcode,
  ADD UNIQUE INDEX idx_product_sku (product_id, sku),
  ADD UNIQUE INDEX idx_product_barcode (product_id, barcode);

ALTER TABLE product_variants
  DROP INDEX idx_product_sku,
  DROP INDEX idx_product_barcode,
  ADD UNIQUE INDEX sku (sku),
  ADD UNIQUE INDEX barcode (barcode);

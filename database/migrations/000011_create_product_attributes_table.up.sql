CREATE TABLE product_attributes (
    product_id BIGINT UNSIGNED NOT NULL,
    attribute_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (product_id, attribute_id),
    CONSTRAINT fk_product_attributes_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_product_attributes_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

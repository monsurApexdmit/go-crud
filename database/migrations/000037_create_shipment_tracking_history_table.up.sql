CREATE TABLE IF NOT EXISTS shipment_tracking_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    shipment_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(50) NOT NULL,
    location VARCHAR(255) NULL,
    description TEXT NULL,
    event_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shipment_id) REFERENCES order_shipments (id) ON DELETE CASCADE,
    INDEX idx_shipment_tracking_shipment_id (shipment_id),
    INDEX idx_shipment_tracking_event_time (event_time)
);

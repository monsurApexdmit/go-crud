CREATE TABLE IF NOT EXISTS attributes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    option_type VARCHAR(50),
    `values` TEXT,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL
);

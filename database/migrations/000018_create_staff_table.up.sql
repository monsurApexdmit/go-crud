CREATE TABLE IF NOT EXISTS staff (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    contact VARCHAR(20),
    joining_date VARCHAR(20),
    role VARCHAR(100),
    status ENUM('Active', 'Inactive') DEFAULT 'Active',
    published TINYINT(1) DEFAULT 0,
    avatar VARCHAR(500),
    salary DECIMAL(12, 2) DEFAULT 0,
    bank_account VARCHAR(100),
    payment_method ENUM('Bank Transfer', 'Cash', 'Check'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_staff_deleted_at ON staff (deleted_at);

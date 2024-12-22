CREATE TABLE transactions (
    id CHAR(10) PRIMARY KEY,
    source_account_number CHAR(12) NOT NULL,
    target_account_number CHAR(12) NOT NULL,
    amount BIGINT NOT NULL,
    bank_id INT NULL,
    type ENUM('internal', 'external', 'debt_payment') NOT NULL,
    description VARCHAR(255) NOT NULL,
    status ENUM('pending', 'success', 'failed') NOT NULL,
    is_source_fee TINYINT(1) NOT NULL,
    source_balance BIGINT NOT NULL,
    target_balance BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);


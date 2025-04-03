CREATE TABLE IF NOT EXISTS accounts (
    address CHAR(42) PRIMARY KEY,
    password_hash TEXT NOT NULL,
    is_validator BOOLEAN DEFAULT FALSE,
    is_zero_address BOOLEAN DEFAULT FALSE,
    is_deposit_address BOOLEAN DEFAULT FALSE,
    sb_balance NUMERIC(20, 8) DEFAULT 0 CHECK (sb_balance >= 0),
    tx_sent_count INTEGER DEFAULT 0 CHECK (tx_sent_count >= 0)
);

CREATE TABLE IF NOT EXISTS stakes (
    id SERIAL PRIMARY KEY,
    account_address CHAR(42) REFERENCES accounts(address),
    amount NUMERIC(20, 8) NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blocks (
    height INTEGER GENERATED ALWAYS AS IDENTITY (START WITH 0 MINVALUE 0) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    size INTEGER NOT NULL,
    parent_hash TEXT,
    hash TEXT NOT NULL,
    fee_recipient_address CHAR(42) NOT NULL REFERENCES accounts(address),
    block_reward NUMERIC(20, 8) NOT NULL CHECK (block_reward > 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_hash CHAR(66) PRIMARY KEY,
    block_height INTEGER REFERENCES blocks(height),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    from_address CHAR(42) REFERENCES accounts(address),
    to_address CHAR(42) REFERENCES accounts(address),
    value NUMERIC(20, 8) NOT NULL CHECK (value > 0),
    transaction_fee NUMERIC(20, 8) NOT NULL CHECK (transaction_fee > 0),
    gas_price NUMERIC(20, 8) NOT NULL CHECK (gas_price > 0),
    gas_used NUMERIC(20, 8) NOT NULL CHECK (gas_used > 0)
);

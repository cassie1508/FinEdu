CREATE TABLE IF NOT EXISTS companies (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    sector VARCHAR(100),
    industry VARCHAR(100),
    market_cap BIGINT,
    revenue BIGINT,
    eps DECIMAL(10,2),
    pe_ratio NUMERIC(10,2),
    dividend_yield NUMERIC(6,4),
    week_high_52 NUMERIC(12,4),
    week_low_52 NUMERIC(12,4),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
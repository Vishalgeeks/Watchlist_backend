CREATE TABLE IF NOT EXISTS portfolios (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id        INT NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    quantity        INT NOT NULL DEFAULT 0,
    avg_price       DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, stock_id)
);

CREATE INDEX IF NOT EXISTS idx_portfolios_user ON portfolios(user_id);
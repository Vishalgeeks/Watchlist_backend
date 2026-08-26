CREATE TABLE IF NOT EXISTS trades (
    id                SERIAL PRIMARY KEY,
    order_id          INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id           INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id          INT NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    side              VARCHAR(10) NOT NULL,
    quantity          INT NOT NULL,
    execution_price   DECIMAL(12,2) NOT NULL,
    total_amount      DECIMAL(18,2) NOT NULL,
    executed_at       TIMESTAMP DEFAULT NOW(),
    CHECK (side IN ('BUY', 'SELL')),
    CHECK (quantity > 0),
    CHECK (execution_price > 0)
);

CREATE INDEX IF NOT EXISTS idx_trades_user ON trades(user_id);
CREATE INDEX IF NOT EXISTS idx_trades_order ON trades(order_id);
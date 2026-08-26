CREATE TABLE IF NOT EXISTS orders (
    id            SERIAL PRIMARY KEY,
    user_id       INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id      INT NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    side          VARCHAR(10) NOT NULL,
    order_type    VARCHAR(20) NOT NULL,
    quantity      INT NOT NULL,
    price         DECIMAL(12,2),
    status        VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CHECK (side IN ('BUY', 'SELL')),
    CHECK (order_type IN ('MARKET', 'LIMIT')),
    CHECK (status IN ('PENDING', 'FILLED', 'CANCELLED', 'REJECTED')),
    CHECK (quantity > 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
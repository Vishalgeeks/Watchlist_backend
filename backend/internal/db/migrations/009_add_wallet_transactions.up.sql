CREATE TABLE IF NOT EXISTS wallet_transactions (
    id              SERIAL PRIMARY KEY,
    wallet_id       INT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(20) NOT NULL,
    amount          DECIMAL(18,2) NOT NULL,
    balance_after   DECIMAL(18,2) NOT NULL,
    reference_type  VARCHAR(50),
    reference_id    INT,
    description     TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    CHECK (type IN ('DEPOSIT', 'WITHDRAWAL', 'BUY_DEBIT', 'SELL_CREDIT')),
    CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_user ON wallet_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_wallet ON wallet_transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_created ON wallet_transactions(created_at DESC);
DROP TABLE IF EXISTS watchlist_items;
DROP TABLE IF EXISTS watchlists;
DROP TABLE IF EXISTS stocks;

-- Purana simple stocks table wapas
CREATE TABLE IF NOT EXISTS stocks (
    id            SERIAL PRIMARY KEY,
    symbol        VARCHAR(20) UNIQUE NOT NULL,
    company_name  VARCHAR(255) NOT NULL,
    exchange      VARCHAR(50),
    ltp DECIMAL(12,2) DEFAULT 0.00,
    last_updated  TIMESTAMP DEFAULT NOW()
);
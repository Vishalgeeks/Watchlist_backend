-- Drop old tables (ulta order)
DROP TABLE IF EXISTS watchlist_items;
DROP TABLE IF EXISTS watchlists;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS users;

-- Recreate with integer IDs
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,        -- auto increment integer
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stocks (
    id            SERIAL PRIMARY KEY,
    symbol        VARCHAR(20) UNIQUE NOT NULL,
    company_name  VARCHAR(255) NOT NULL,
    exchange      VARCHAR(50),
    LTP DECIMAL(12,2) DEFAULT 0.00,
    last_updated  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS watchlists (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS watchlist_items (
    id           SERIAL PRIMARY KEY,
    watchlist_id INT NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    stock_id     INT NOT NULL REFERENCES stocks(id)     ON DELETE CASCADE,
    added_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(watchlist_id, stock_id)
);

-- Sample Data
INSERT INTO users (name) VALUES ('Vish'), ('Rahul');

INSERT INTO stocks (symbol, company_name, exchange, LTP) VALUES
    ('RELIANCE',   'Reliance Industries Ltd',   'NSE', 2450.75),
    ('TCS',        'Tata Consultancy Services', 'NSE', 3820.50),
    ('INFY',       'Infosys Ltd',               'NSE', 1456.30),
    ('HDFCBANK',   'HDFC Bank Ltd',             'NSE', 1678.90),
    ('TATAMOTORS', 'Tata Motors Ltd',           'NSE',  875.40),
    ('WIPRO',      'Wipro Ltd',                 'NSE',  456.20),
    ('ICICIBANK',  'ICICI Bank Ltd',            'NSE', 1123.60),
    ('SBIN',       'State Bank of India',       'NSE',  734.85);
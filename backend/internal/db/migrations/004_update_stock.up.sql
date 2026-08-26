-- Purani stocks table drop karo
DROP TABLE IF EXISTS watchlist_items;
DROP TABLE IF EXISTS watchlists;
DROP TABLE IF EXISTS stocks;

-- Naya stocks table with all CSV columns
CREATE TABLE IF NOT EXISTS stocks (
    id                      SERIAL PRIMARY KEY,
    exchange_instrument_id  VARCHAR(50),
    segment                 VARCHAR(50),
    instrument_type         VARCHAR(50),
    symbol                  VARCHAR(50) UNIQUE NOT NULL,
    display_name            VARCHAR(255),
    company_name            VARCHAR(255),
    isin                    VARCHAR(50),
    series                  VARCHAR(20),
    exchange                VARCHAR(50),
    contract_expiration     VARCHAR(50),
    strike                  DECIMAL(12,2),
    option_type             VARCHAR(10),
    underlying_symbol_id    VARCHAR(50),
    underlying_symbol       VARCHAR(50),
    lot_size                INT,
    tick_size               DECIMAL(12,4),
    upper_circuit           DECIMAL(12,2),
    lower_circuit           DECIMAL(12,2),
    freeze_qty              INT,
    description             TEXT,
    ltp                     DECIMAL(12,2),
    open                    DECIMAL(12,2),
    high                    DECIMAL(12,2),
    low                     DECIMAL(12,2),
    close                   DECIMAL(12,2),
    vol                     BIGINT,
    oi                      BIGINT,
    bid                     DECIMAL(12,2),
    ask                     DECIMAL(12,2),
    bid_qty                 INT,
    ask_qty                 INT,
    cautionary_message_info TEXT,
    last_updated            TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stocks_symbol ON stocks(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_exchange ON stocks(exchange);

-- Watchlist tables recreate
CREATE TABLE IF NOT EXISTS watchlists (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_watchlists_user ON watchlists(user_id);

CREATE TABLE IF NOT EXISTS watchlist_items (
    id           SERIAL PRIMARY KEY,
    watchlist_id INT NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    stock_id     INT NOT NULL REFERENCES stocks(id)     ON DELETE CASCADE,
    added_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(watchlist_id, stock_id)
);

CREATE INDEX IF NOT EXISTS idx_items_watchlist ON watchlist_items(watchlist_id);
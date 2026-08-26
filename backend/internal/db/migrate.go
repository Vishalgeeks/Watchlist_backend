package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(connStr string) error {

	m, err := migrate.New(
		"file://./internal/db/migrations",
		connStr,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {

		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to run")
			return ensureMissingTables(connStr)
		}

		if strings.Contains(err.Error(), "Dirty database version") {
			log.Println("Dirty database detected, attempting repair")
			return repairMigrations(connStr)
		}

		return err
	}

	return ensureMissingTables(connStr)
}

func repairMigrations(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	var dirtyVersion int
	err = db.QueryRow(`SELECT version FROM schema_migrations WHERE dirty = true LIMIT 1`).Scan(&dirtyVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No dirty version found in schema_migrations")
			return nil
		}
		return err
	}

	log.Printf("Dirty version: %d", dirtyVersion)

	forceVersion := 0
	if dirtyVersion > 1 {
		forceVersion = dirtyVersion - 1
	}
	log.Printf("Forcing version to %d", forceVersion)

	m, err := migrate.New(
		"file://./internal/db/migrations",
		connStr,
	)
	if err != nil {
		return err
	}

	m.Force(forceVersion)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up after force failed: %w", err)
	}

	return ensureMissingTables(connStr)
}

func ensureMissingTables(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	missing, err := getMissingTables(db)
	if err != nil {
		return err
	}

	if len(missing) == 0 {
		return nil
	}

	log.Printf("Missing tables detected: %v", missing)

	for _, table := range missing {
		sqlStmt := getCreateTableSQL(table)
		if sqlStmt == "" {
			continue
		}
		if _, err := db.Exec(sqlStmt); err != nil {
			log.Printf("Failed to create table %s: %v", table, err)
		} else {
			log.Printf("Created missing table: %s", table)
		}
	}

	return nil
}

func getMissingTables(db *sql.DB) ([]string, error) {
	expected := []string{"orders", "trades"}
	var missing []string
	for _, table := range expected {
		var exists bool
		err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)`, table).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			missing = append(missing, table)
		}
	}
	return missing, nil
}

func getCreateTableSQL(table string) string {
	switch table {
	case "orders":
		return `
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
		`
	case "trades":
		return `
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
		`
	default:
		return ""
	}
}

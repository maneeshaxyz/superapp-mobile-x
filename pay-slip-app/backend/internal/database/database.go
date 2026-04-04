// internal/database/database.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pay-slip-app/internal/configs"
	"pay-slip-app/internal/constants"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Conn *sql.DB
	mu   sync.RWMutex
}

// NewDatabase creates a new database connection with pool tuning and a background health pinger.
func NewDatabase(cfg configs.DBConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(constants.ConnMaxLifetimeMinutes) * time.Minute)
	db.SetMaxIdleConns(constants.MaxIdleConns)
	db.SetMaxOpenConns(constants.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	d := &Database{Conn: db}
	log.Println("Database connection established")

	// Background pinger — keeps connections alive, auto-reconnects on repeated failures.
	go func(dsn string, database *Database) {
		ticker := time.NewTicker(time.Duration(constants.PingIntervalSeconds) * time.Second)
		defer ticker.Stop()
		failCount := 0
		for range ticker.C {
			err := database.Ping()
			if err != nil {
				log.Printf("DB ping failed: %v", err)
				failCount++
			} else {
				failCount = 0
				continue
			}

			if failCount >= constants.ReconnectFailThreshold {
				log.Println("Attempting DB reconnect after repeated ping failures")
				newDB, err := sql.Open("mysql", dsn)
				if err != nil {
					log.Printf("reconnect: sql.Open error: %v", err)
					continue
				}
				newDB.SetConnMaxLifetime(time.Duration(constants.ConnMaxLifetimeMinutes) * time.Minute)
				newDB.SetMaxIdleConns(constants.MaxIdleConns)
				newDB.SetMaxOpenConns(constants.MaxOpenConns)
				if err := newDB.Ping(); err != nil {
					log.Printf("reconnect: ping failed: %v", err)
					_ = newDB.Close()
					continue
				}

				database.mu.Lock()
				old := database.Conn
				database.Conn = newDB
				database.mu.Unlock()
				_ = old.Close()
				log.Println("DB reconnect successful")
				failCount = 0
			}
		}
	}(dsn, d)

	return d, nil
}

// ── Database Methods ───────────────────────────────────────────────────────

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.Exec(query, args...)
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.Query(query, args...)
}

func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.Conn.BeginTx(ctx, opts)
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.QueryRow(query, args...)
}

func (db *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.Conn.ExecContext(ctx, query, args...)
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.Conn.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.Conn.QueryRowContext(ctx, query, args...)
}

func (db *Database) Ping() error {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.Ping()
}

func (db *Database) Close() error {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.Close()
}

func (db *Database) Begin() (*sql.Tx, error) {
	db.mu.RLock()
	conn := db.Conn
	db.mu.RUnlock()
	return conn.Begin()
}

package database

import (
	"context"
	"database/sql"
	"fmt"
	"stock-analyser/utils"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitializeConnection(ctx context.Context, config *utils.AppConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	DB = db
	return nil
}

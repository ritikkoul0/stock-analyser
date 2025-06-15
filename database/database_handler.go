package database

import (
	"context"
	"fmt"
	"log"
)

func SaveUser(ctx context.Context, username, email, hashedPassword string) error {
	// Attempt to insert the user
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
`

	var userID int
	err := DB.QueryRowContext(ctx, query, username, email, hashedPassword).Scan(&userID)
	if err != nil {
		log.Printf("❌ Failed to insert user: %v", err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("✅ User inserted with ID: %d\n", userID)
	return nil
}

func UserExists(ctx context.Context, username string, email string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1 OR email = $2
		);
	`
	err := DB.QueryRow(query, username, email).Scan(&exists)
	return exists, err
}

func GetStock(ctx context.Context) ([]string, error) {
	query := `SELECT stock_symbol FROM stocksymbol;`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []string
	for rows.Next() {
		var stock string
		if err := rows.Scan(&stock); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

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

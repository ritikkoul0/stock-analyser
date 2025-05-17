package database

import (
	"context"
	"fmt"
	"log"
)

func SaveUser(ctx context.Context, username, email, hashedPassword string) error {
	// Diagnostic: list all tables in 'public' schema
	rows, err := DB.QueryContext(ctx, "SELECT tablename FROM pg_tables WHERE schemaname='public'")
	if err != nil {
		log.Printf("❌ Failed to query table names: %v", err)
	} else {
		defer rows.Close()
		log.Println("📋 Available tables in 'public' schema:")
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				log.Printf(" - %s", tableName)
			}
		}
	}

	// Attempt to insert the user
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
`

	var userID int
	err = DB.QueryRowContext(ctx, query, username, email, hashedPassword).Scan(&userID)
	if err != nil {
		log.Printf("❌ Failed to insert user: %v", err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("✅ User inserted with ID: %d\n", userID)
	return nil
}

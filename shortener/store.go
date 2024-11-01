package shortener

import (
	"context"
	"fmt"
	"time"
)

// SaveURL stores the original URL, short code, and expiration time in the database.
func SaveURL(originalURL, shortCode string, expiresAt *time.Time) error {
	query := `
        INSERT INTO urls (original_url, short_code, expires_at)
        VALUES ($1, $2, $3)
    `

	_, err := DB.Exec(context.Background(), query, originalURL, shortCode, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save URL: %v", err)
	}
	return nil
}

// GetOriginalURL retrieves the original URL by its short code.
func GetOriginalURL(shortCode string) (string, error) {
	var originalURL string
	var expiresAt *time.Time

	query := `
        SELECT original_url, expires_at 
        FROM urls 
        WHERE short_code = $1
    `

	err := DB.QueryRow(context.Background(), query, shortCode).Scan(&originalURL, &expiresAt)
	if err != nil {
		return "", fmt.Errorf("URL not found: %v", err)
	}

	// Check if the URL has expired
	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return "", fmt.Errorf("URL has expired")
	}

	return originalURL, nil
}

// CheckCodeExists checks if a custom code already exists in the database.
func CheckCodeExists(code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`

	err := DB.QueryRow(context.Background(), query, code).Scan(&exists)
	return exists, err
}

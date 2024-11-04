package shortener

import (
	"context"
	"fmt"
	"time"
)

// SaveURL stores the original URL, short code, and expiration time in the database.
func SaveURL(originalURL, name string, shortCode string, customCode bool, expiresAt *time.Time, customExpiry bool, created_by int) error {
	query := `
        INSERT INTO urls (original_url, name, short_code, is_custom_code, expires_at, is_custom_expiry, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := DB.Exec(context.Background(), query, originalURL, name, shortCode, customCode, expiresAt, customExpiry, created_by)
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

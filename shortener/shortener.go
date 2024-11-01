// shortener/shortener.go
package shortener

import (
	"encoding/json"
	"net/http"
	"time"
)

// URLRequest represents the structure of the incoming JSON request.
type URLRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
	ExpiresAt  string `json:"expires_at,omitempty"`
}

// URLResponse represents the JSON structure of the response.
type URLResponse struct {
	OriginalURL string `json:"original_url"`
	Code        string `json:"code"`
}

// ShortenURL handles POST requests to shorten URLs.
func ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid request body, 'url' field is required", http.StatusBadRequest)
		return
	}

	// Handle custom code if provided
	code := req.CustomCode
	if code == "" {
		code = GenerateShortCode(req.URL, 8) // Generate a unique code if not provided
	} else {
		// Check if the custom code is already in use
		exists, err := CheckCodeExists(code)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Custom code already in use", http.StatusConflict)
			return
		}
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			http.Error(w, "Invalid 'expires_at' format, expected RFC3339 format", http.StatusBadRequest)
			return
		}
		expiresAt = &parsedTime
	}

	// Store the URL with code and optional expiration date
	err = SaveURL(req.URL, code, expiresAt)
	if err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	response := URLResponse{
		OriginalURL: req.URL,
		Code:        code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

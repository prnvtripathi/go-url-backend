package shortener

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type URLRequest struct {
	URL               string `json:"url"`
	Name              string `json:"name"`
	CustomCode        string `json:"custom_code,omitempty"`
	CustomCodeEnabled bool   `json:"custom_code_enabled,omitempty"`
	ExpiresAt         string `json:"expires_at"`
	CustomExpiry      bool   `json:"custom_expiry,omitempty"`
	CreatedBy         int    `json:"created_by"`
}

// URLResponse represents the JSON structure of the response.
type URLResponse struct {
	Success     bool   `json:"success"`
	OriginalURL string `json:"original_url"`
	Code        string `json:"code"`
	Message     string `json:"message,omitempty"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		response := URLResponse{
			Success: false,
			Message: "Invalid request",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	code := req.CustomCode
	if code == "" {
		code = GenerateShortCode(req.URL, 8)
	} else {
		exists, err := CheckCodeExists(code)
		if err != nil {
			log.Printf("Error checking custom code: %v", err)
			response := URLResponse{
				Success: false,
				Message: "Failed to check custom code",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if exists {
			response := URLResponse{
				Success: false,
				Message: "Custom code already exists",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			response := URLResponse{
				Success: false,
				Message: "Invalid 'expires_at' format, expected RFC3339 format",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		expiresAt = &parsedTime
	}

	err = SaveURL(req.URL, req.Name, code, req.CustomCodeEnabled, expiresAt, req.CustomExpiry, req.CreatedBy)
	if err != nil {
		log.Printf("Failed to store URL: %v", err)
		response := URLResponse{
			Success: false,
			Message: "Failed to store URL",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := URLResponse{
		Success:     true,
		OriginalURL: req.URL,
		Code:        code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

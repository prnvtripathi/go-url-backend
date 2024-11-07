package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prnvtripathi/go-url-api/shortener"
)

// URLRequest represents the request payload to get URLs.
type URLRequest struct {
	UserID int `json:"user_id"`
}

// URLResponse represents the response payload.
type URLResponse struct {
	Success bool            `json:"success"`
	URLs    []shortener.URL `json:"urls,omitempty"`
	Message string          `json:"message,omitempty"`
}

// getUrlsHandler handles the HTTP request to retrieve all URLs for a user.
func getUrlsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request into URLRequest
	var req URLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.UserID == 0 {
		response := URLResponse{
			Success: false,
			Message: "Invalid request",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Printf("Invalid request: %v", err)
		return
	}

	// Call GetAllUrls from the shortener package
	urls, err := shortener.GetAllUrls(req.UserID)
	if err != nil {
		response := URLResponse{
			Success: false,
			Message: "Failed to retrieve URLs",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Printf("Failed to retrieve URLs: %v", err)
		return
	}

	// Prepare and send response
	response := URLResponse{
		Success: true,
		URLs:    urls,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

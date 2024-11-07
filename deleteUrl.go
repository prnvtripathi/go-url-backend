package main

import (
	"encoding/json"
	"net/http"

	"github.com/prnvtripathi/go-url-api/shortener"
)

type DeleteUrlRequest struct {
	UrlId  int `json:"url_id"`
	UserId int `json:"user_id"`
}

type DeleteUrlResponse struct {
	Sucess  bool   `json:"success"`
	Message string `json:"message"`
}

func deleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request DeleteUrlRequest
	var response DeleteUrlResponse

	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Message = "Invalid request"
		json.NewEncoder(w).Encode(response)
		return
	}

	err = shortener.DeleteUrl(request.UrlId, request.UserId)
	if err != nil {
		response.Message = "Failed to delete URL"
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Sucess = true
	response.Message = "URL deleted successfully"
	json.NewEncoder(w).Encode(response)
}

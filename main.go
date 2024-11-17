package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/prnvtripathi/go-url-api/redirect"
	"github.com/prnvtripathi/go-url-api/shortener"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running")
}

// main is the entry point of the application. It sets up an HTTP server
// that listens on port 8080 and handles requests to the "/health" endpoint
// using the healthCheck handler function. If the server fails to start,
// it logs the error and exits.
func main() {

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env.local file")
	}

	// Connect to the database
	err = shortener.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer shortener.CloseDB()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is up and running")
	})
	http.HandleFunc("/shorten", shortener.ShortenURL) // Shorten URL handler
	http.HandleFunc("/r/", redirect.RedirectHandler)  // Redirect handler
	http.HandleFunc("/getUrls", getUrlsHandler)       // Get URLs handler
	http.HandleFunc("/deleteUrl", deleteUrlHandler)   // Delete URL handler

	// Start the server
	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

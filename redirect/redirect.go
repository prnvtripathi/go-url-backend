// redirect/redirect.go
package redirect

import (
	"net/http"
	"strings"

	"github.com/prnvtripathi/go-url-api/shortener"
)

// RedirectHandler handles the redirection based on the short code in the URL.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the code from the URL path, expected as "/r/{code}"
	code := strings.TrimSpace(r.URL.Path[len("/r/"):])

	if code == "" {
		http.Error(w, "Code not provided", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL associated with the code from the database
	originalURL, err := shortener.GetOriginalURL(code)
	if err != nil {
		http.Error(w, "URL not found or has expired", http.StatusNotFound)
		return
	}

	// Ensure the URL is fully qualified
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

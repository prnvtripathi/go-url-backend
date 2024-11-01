// shortener/generator.go
package shortener

import (
	"crypto/sha256"
	"encoding/binary"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode creates a unique base-62 short code for a given URL.
func GenerateShortCode(url string, length int) string {
	// Step 1: Hash the URL
	hash := sha256.Sum256([]byte(url))

	// Step 2: Take the first 8 bytes of the hash as a number
	hashPrefix := binary.BigEndian.Uint64(hash[:8])

	// Step 3: Convert this number to a base-62 encoded string
	return base62Encode(hashPrefix, length)
}

// base62Encode converts an integer to a base-62 encoded string of fixed length.
func base62Encode(number uint64, length int) string {
	encoded := make([]byte, length)
	base := uint64(len(base62Chars))

	// Fill the encoded string backwards to ensure a fixed length
	for i := length - 1; i >= 0; i-- {
		encoded[i] = base62Chars[number%base]
		number = number / base
	}
	return string(encoded)
}

// Package crypto provides cryptographic utility functions.
package crypto

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

// SHA1 generates the SHA-1 hash of a given string.
// It takes a string as input and returns its 40-character hexadecimal SHA-1 hash.
func SHA1(s string) string {
	// Convert the input string to a byte slice.
	b := []byte(s)
	// Compute the SHA-1 hash of the byte slice.
	h := sha1.Sum(b)
	// Format the hash as a hexadecimal string and return it.
	return fmt.Sprintf("%x", h)
}

// SHA256 generates the SHA-256 hash of a given string.
// It takes a string as input and returns its 64-character hexadecimal SHA-256 hash.
func SHA256(s string) string {
	// Convert the input string to a byte slice.
	b := []byte(s)
	// Compute the SHA-256 hash of the byte slice.
	h := sha256.Sum256(b)
	// Format the hash as a hexadecimal string and return it.
	return fmt.Sprintf("%x", h)
}

// SHA512 generates the SHA-512 hash of a given string.
// It takes a string as input and returns its 128-character hexadecimal SHA-512 hash.
func SHA512(s string) string {
	// Convert the input string to a byte slice.
	b := []byte(s)
	// Compute the SHA-512 hash of the byte slice.
	h := sha512.Sum512(b)
	// Format the hash as a hexadecimal string and return it.
	return fmt.Sprintf("%x", h)
}

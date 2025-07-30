// Package crypto provides cryptographic utility functions.
package crypto

import (
	"crypto/md5"
	"fmt"
)

// MD5 generates the MD5 hash of a given string.
// It takes a string as input and returns its 32-character hexadecimal MD5 hash.
func MD5(s string) string {
	// Convert the input string to a byte slice.
	b := []byte(s)
	// Compute the MD5 hash of the byte slice.
	h := md5.Sum(b)
	// Format the hash as a hexadecimal string and return it.
	return fmt.Sprintf("%x", h)
}

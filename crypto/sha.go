package crypto

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func SHA1(s string) string {
	b := []byte(s)
	h := sha1.Sum(b)
	return fmt.Sprintf("%x", h)
}

func SHA256(s string) string {
	b := []byte(s)
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h)
}

func SHA512(s string) string {
	b := []byte(s)
	h := sha512.Sum512(b)
	return fmt.Sprintf("%x", h)
}

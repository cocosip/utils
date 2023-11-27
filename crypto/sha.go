package crypto

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

func SHA1(s string) string {
	b := []byte(s)
	h := sha1.Sum(b)
	return fmt.Sprintf("%x", h)
}

func SHA256(s string) string {
	b := []byte(s)
	h := sha256.New()
	return fmt.Sprintf("%x", h.Sum(b))
}

package crypto

import (
	"crypto/md5"
	"fmt"
)

func MD5(s string) string {
	b := []byte(s)
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}

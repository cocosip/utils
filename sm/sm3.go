// Package sm provides SM2, SM3, and SM4 cryptographic functionalities.
// It leverages the 'github.com/tjfoc/gmsm' library.
package sm

import (
	"fmt"

	"github.com/tjfoc/gmsm/sm3"
)

// HashSM3 computes the SM3 hash of the input data.
// `data` is the input data to be hashed (raw bytes).
// It returns the hex-encoded hash value.
func HashSM3(data []byte) (hashHex string, err error) {
	h := sm3.New()
	_, err = h.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write data to SM3 hasher: %w", err)
	}
	hashBytes := h.Sum(nil)
	return fmt.Sprintf("%x", hashBytes), nil
}

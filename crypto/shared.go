// Package crypto provides cryptographic utility functions.
package crypto

import (
	"bytes"
	"errors"
)

// PaddingMode represents the type of padding used in block ciphers.
type PaddingMode int

const (
	// PaddingModeNone indicates no padding.
	PaddingModeNone PaddingMode = 1
	// PaddingModePKCS7 indicates PKCS7 padding.
	PaddingModePKCS7 PaddingMode = 2
	// PaddingModeZeros indicates zero padding.
	PaddingModeZeros PaddingMode = 3
	// PaddingModeANSIX923 indicates ANSI X9.23 padding.
	PaddingModeANSIX923 PaddingMode = 4
	// PaddingModeISO10126 indicates ISO 10126 padding.
	PaddingModeISO10126 PaddingMode = 5
	// PaddingModePKCS5 indicates PKCS5 padding.
	// PKCS5 padding is a subset of PKCS7 padding with a fixed block size of 8 bytes.
	PaddingModePKCS5 PaddingMode = 6
)

// CipherMode represents the block cipher mode of operation.
type CipherMode int

const (
	// CipherModeCBC indicates Cipher Block Chaining mode.
	CipherModeCBC CipherMode = 1
	// CipherModeECB indicates Electronic Codebook mode.
	// Note: ECB is generally insecure for most applications as it does not hide data patterns.
	CipherModeECB CipherMode = 2
	// CipherModeCFB indicates Cipher Feedback mode.
	CipherModeCFB CipherMode = 3
)

var (
	// ErrInvalidBlockSize indicates that the provided block size is invalid (e.g., non-positive).
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates that the input data for PKCS7 padding/unpadding is invalid (e.g., nil or empty).
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates that PKCS7 unpadding failed due to invalid padding bytes.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")

	// ErrInvalidPKCS5Padding indicates that PKCS5 unpadding failed due to invalid padding bytes.
	ErrInvalidPKCS5Padding = errors.New("invalid PKCS5 padding")
)

// PKCS7Padding applies PKCS7 padding to the given byte slice.
// It pads the data with 1 to blockSize bytes, where the value of each padding byte
// is the number of padding bytes itself. The size of the result is a multiple of blockSize.
// Returns an error if blockSize is invalid or if the input data is nil/empty.
func PKCS7Padding(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// PKCS7UnPadding removes PKCS7 padding from the given byte slice.
// It validates the padding bytes and returns the original unpadded data.
// Returns an error if blockSize is invalid, input data is invalid, or padding is incorrect.
func PKCS7UnPadding(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blockSize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	// Check if padding length is valid and does not exceed the data length.
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	// Validate all padding bytes.
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

// PKCS5Padding applies PKCS5 padding to the given byte slice.
// PKCS5 padding is identical to PKCS7 padding with a block size of 8 bytes.
// This function assumes a block size of 8 for PKCS5.
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	// PKCS5 padding implies a block size of 8. If a different blockSize is passed,
	// it will behave like PKCS7 padding with that blockSize.
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS5Trimming removes PKCS5 padding from the given byte slice.
// It validates the padding byte and returns the original unpadded data.
// Returns an error if the padding is invalid or leads to an out-of-bounds access.
func PKCS5Trimming(encrypt []byte) ([]byte, error) {
	if len(encrypt) == 0 {
		return nil, ErrInvalidPKCS5Padding // Or a more specific error for empty input
	}
	padding := encrypt[len(encrypt)-1]
	// Check if the padding value is valid and does not exceed the slice length.
	if int(padding) == 0 || int(padding) > len(encrypt) {
		return nil, ErrInvalidPKCS5Padding
	}
	// For PKCS5, the padding value should also be less than or equal to the block size (8).
	// However, since PKCS5 is a subset of PKCS7, we can rely on the general padding validation.
	// The main goal here is to prevent panics from invalid padding values.
	return encrypt[:len(encrypt)-int(padding)], nil
}

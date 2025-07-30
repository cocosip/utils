// Package crypto provides cryptographic utility functions.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var (
	// DefaultAESKey is a default 16-byte AES key (AES-128).
	DefaultAESKey = []byte{0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00}
	// DefaultAESIV is a default 16-byte AES Initialization Vector (IV).
	DefaultAESIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
)

// AESCrypto represents an AES encryption/decryption context.
// It holds the key, IV, padding mode, and cipher mode.
type AESCrypto struct {
	key         []byte
	iv          []byte
	paddingMode PaddingMode
	cipherMode  CipherMode
}

// NewAES creates a new AESCrypto instance with default settings.
// Default settings are AES-128, CBC mode, and PKCS7 padding.
func NewAES() *AESCrypto {
	return &AESCrypto{
		key:         DefaultAESKey,
		iv:          DefaultAESIV,
		paddingMode: PaddingModePKCS7,
		cipherMode:  CipherModeCBC,
	}
}

// WithKey sets the AES key for the AESCrypto instance.
// The key length must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256 respectively.
func (c *AESCrypto) WithKey(key []byte) *AESCrypto {
	c.key = key
	return c
}

// WithIV sets the AES Initialization Vector (IV) for the AESCrypto instance.
// The IV length must be equal to the cipher block size (aes.BlockSize, typically 16 bytes).
func (c *AESCrypto) WithIV(iv []byte) *AESCrypto {
	c.iv = iv
	return c
}

// WithMode sets the padding mode and cipher mode for the AESCrypto instance.
func (c *AESCrypto) WithMode(paddingMode PaddingMode, cipherMode CipherMode) *AESCrypto {
	c.paddingMode = paddingMode
	c.cipherMode = cipherMode
	return c
}

// EncryptToBase64 encrypts the given plaintext string and returns the result as a base64-encoded string.
func (c *AESCrypto) EncryptToBase64(plaintext string) (string, error) {
	cipherBuffer, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt plaintext: %w", err)
	}
	return base64.StdEncoding.EncodeToString(cipherBuffer), nil
}

// DecryptFromBase64 decrypts a base64-encoded ciphertext string and returns the original plaintext string.
func (c *AESCrypto) DecryptFromBase64(ciphertext string) (string, error) {
	cipherBuffer, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	plainBuffer, err := c.Decrypt(cipherBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}
	return string(plainBuffer), nil
}

// EncryptPlainText encrypts the given plaintext string and returns the result as a hex-encoded string.
func (c *AESCrypto) EncryptPlainText(plaintext string) (string, error) {
	cipherBuffer, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt plaintext: %w", err)
	}
	return hex.EncodeToString(cipherBuffer), nil
}

// DecryptCipherText decrypts a hex-encoded ciphertext string and returns the original plaintext string.
func (c *AESCrypto) DecryptCipherText(ciphertext string) (string, error) {
	cipherBuffer, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex ciphertext: %w", err)
	}

	plainBuffer, err := c.Decrypt(cipherBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}
	return string(plainBuffer), nil
}

// Encrypt encrypts the given plaintext byte slice using the configured AES key, IV, padding, and cipher mode.
// It returns the ciphertext as a byte slice.
func (c *AESCrypto) Encrypt(plainBuffer []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	blockSize := block.BlockSize()
	if len(c.iv) != blockSize {
		return nil, fmt.Errorf("IV length (%d) must be equal to block size (%d)", len(c.iv), blockSize)
	}

	switch c.paddingMode {
	case PaddingModePKCS7:
		plainBuffer, err = PKCS7Padding(plainBuffer, blockSize)
		if err != nil {
			return nil, fmt.Errorf("PKCS7 padding failed: %w", err)
		}
	case PaddingModePKCS5:
		plainBuffer = PKCS5Padding(plainBuffer, blockSize)
	default:
		return nil, fmt.Errorf("unsupported padding mode: %v", c.paddingMode)
	}

	ciphertext := make([]byte, len(plainBuffer))
	var bm cipher.BlockMode
	switch c.cipherMode {
	case CipherModeCBC:
		bm = cipher.NewCBCEncrypter(block, c.iv)
	case CipherModeECB:
		// ECB mode is generally not recommended for most use cases due to lack of diffusion.
		// If you need ECB, you'd implement it here.
		return nil, fmt.Errorf("ECB cipher mode not implemented")
	default:
		return nil, fmt.Errorf("unsupported cipher mode: %v", c.cipherMode)
	}

	bm.CryptBlocks(ciphertext, plainBuffer)
	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext byte slice using the configured AES key, IV, padding, and cipher mode.
// It returns the plaintext as a byte slice.
func (c *AESCrypto) Decrypt(cipherBuffer []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	blockSize := block.BlockSize()
	if len(c.iv) != blockSize {
		return nil, fmt.Errorf("IV length (%d) must be equal to block size (%d)", len(c.iv), blockSize)
	}

	if len(cipherBuffer)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(cipherBuffer))
	var bm cipher.BlockMode
	switch c.cipherMode {
	case CipherModeCBC:
		bm = cipher.NewCBCDecrypter(block, c.iv)
	case CipherModeECB:
		// ECB mode is generally not recommended for most use cases due to lack of diffusion.
		// If you need ECB, you'd implement it here.
		return nil, fmt.Errorf("ECB cipher mode not implemented")
	default:
		return nil, fmt.Errorf("unsupported cipher mode: %v", c.cipherMode)
	}

	bm.CryptBlocks(plaintext, cipherBuffer)

	switch c.paddingMode {
	case PaddingModePKCS7:
		var err error
		plaintext, err = PKCS7UnPadding(plaintext, blockSize)
		if err != nil {
			return nil, fmt.Errorf("PKCS7 unpadding failed: %w", err)
		}
	case PaddingModePKCS5:
		plaintext, err = PKCS5Trimming(plaintext)
		if err != nil {
			return nil, fmt.Errorf("PKCS5 trimming failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported padding mode: %v", c.paddingMode)
	}

	return plaintext, nil
}

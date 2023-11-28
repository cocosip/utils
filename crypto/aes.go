package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var (
	DefaultAESKey = []byte{0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00}
	DefaultAESIV  = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
)

type AESCrypto struct {
	key         []byte
	iv          []byte
	paddingMode PaddingMode
	cipherMode  CipherMode
}

func NewAES() *AESCrypto {
	return &AESCrypto{
		key:         DefaultAESKey,
		iv:          DefaultAESIV,
		paddingMode: PaddingModePKCS7,
		cipherMode:  CipherModeCBC,
	}
}

// WithKey , 16, 24, 32.  AES-128, AES-192, AES-256.
func (c *AESCrypto) WithKey(key []byte) (*AESCrypto, error) {
	c.key = key
	return c, nil
}

func (c *AESCrypto) WithIV(iv []byte) *AESCrypto {
	c.iv = iv
	return c
}

func (c *AESCrypto) WithMode(paddingMode PaddingMode, cipherMode CipherMode) *AESCrypto {
	c.paddingMode = paddingMode
	c.cipherMode = cipherMode
	return c
}

func (c *AESCrypto) EncryptToBase64(plaintext string) (string, error) {
	cipherBuffer, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBuffer), nil
}

func (c *AESCrypto) DecryptFromBase64(ciphertext string) (string, error) {
	cipherBuffer, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plainBuffer, err := c.Decrypt(cipherBuffer)
	if err != nil {
		return "", err
	}
	return string(plainBuffer), nil
}

func (c *AESCrypto) EncryptPlainText(plaintext string) (string, error) {
	cipherBuffer, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipherBuffer), nil
}

func (c *AESCrypto) DecryptCipherText(ciphertext string) (string, error) {
	cipherBuffer, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plainBuffer, err := c.Decrypt(cipherBuffer)
	if err != nil {
		return "", err
	}
	return string(plainBuffer), nil
}

func (c *AESCrypto) Encrypt(plainBuffer []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	switch c.paddingMode {
	case PaddingModePKCS7:
		plainBuffer, err = PKCS7Padding(plainBuffer, block.BlockSize())
		if err != nil {
			return nil, err
		}
		break
	case PaddingModePKCS5:
		plainBuffer = PKCS5Padding(plainBuffer, block.BlockSize())
		break
	default:
		return nil, fmt.Errorf("invalid padding %v", c.paddingMode)
	}

	cipherBuffer := make([]byte, block.BlockSize()+len(plainBuffer))
	var bm cipher.BlockMode
	switch c.cipherMode {
	case CipherModeCBC:
		bm = cipher.NewCBCEncrypter(block, c.iv)
		break
	default:
		return nil, fmt.Errorf("invalid cipher %v", c.cipherMode)
	}
	bm.CryptBlocks(cipherBuffer[block.BlockSize():], cipherBuffer)
	return cipherBuffer, nil
}

func (c *AESCrypto) Decrypt(cipherBuffer []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	cipherBuffer = cipherBuffer[block.BlockSize():]
	var bm cipher.BlockMode
	switch c.cipherMode {
	case CipherModeCBC:
		bm = cipher.NewCBCEncrypter(block, c.iv)
		break
	default:
		return nil, fmt.Errorf("invalid cipher %v", c.cipherMode)
	}
	bm.CryptBlocks(cipherBuffer, cipherBuffer)
	switch c.paddingMode {
	case PaddingModePKCS7:
		cipherBuffer, err = PKCS7UnPadding(cipherBuffer, block.BlockSize())
		if err != nil {
			return nil, err
		}
		break
	case PaddingModePKCS5:
		cipherBuffer = PKCS5Trimming(cipherBuffer)
		break
	default:
		return nil, fmt.Errorf("invalid padding %v", c.paddingMode)
	}

	return cipherBuffer, nil
}

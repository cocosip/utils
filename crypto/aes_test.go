package crypto

import (
	"bytes"
	"testing"
)

func TestNewAES(t *testing.T) {
	aes := NewAES()
	if aes == nil {
		t.Fatal("NewAES returned nil")
	}
	if !bytes.Equal(aes.key, DefaultAESKey) {
		t.Errorf("NewAES default key mismatch, got %v, want %v", aes.key, DefaultAESKey)
	}
	if !bytes.Equal(aes.iv, DefaultAESIV) {
		t.Errorf("NewAES default IV mismatch, got %v, want %v", aes.iv, DefaultAESIV)
	}
	if aes.paddingMode != PaddingModePKCS7 {
		t.Errorf("NewAES default padding mode mismatch, got %v, want %v", aes.paddingMode, PaddingModePKCS7)
	}
	if aes.cipherMode != CipherModeCBC {
		t.Errorf("NewAES default cipher mode mismatch, got %v, want %v", aes.cipherMode, CipherModeCBC)
	}
}

func TestAESCrypto_WithKey(t *testing.T) {
	aes := NewAES()
	newKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	aes.WithKey(newKey)
	if !bytes.Equal(aes.key, newKey) {
		t.Errorf("WithKey failed, got %v, want %v", aes.key, newKey)
	}
}

func TestAESCrypto_WithIV(t *testing.T) {
	aes := NewAES()
	newIV := []byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	aes.WithIV(newIV)
	if !bytes.Equal(aes.iv, newIV) {
		t.Errorf("WithIV failed, got %v, want %v", aes.iv, newIV)
	}
}

func TestAESCrypto_WithMode(t *testing.T) {
	aes := NewAES()
	aes.WithMode(PaddingModePKCS5, CipherModeECB)
	if aes.paddingMode != PaddingModePKCS5 {
		t.Errorf("WithMode padding mode failed, got %v, want %v", aes.paddingMode, PaddingModePKCS5)
	}
	if aes.cipherMode != CipherModeECB {
		t.Errorf("WithMode cipher mode failed, got %v, want %v", aes.cipherMode, CipherModeECB)
	}
}

func TestAESCrypto_EncryptDecrypt(t *testing.T) {
	tests := []struct {
		name        string
		plaintext   string
		key         []byte
		iv          []byte
		paddingMode PaddingMode
		cipherMode  CipherMode
		expectError bool
	}{
		{
			name:        "Default_CBC_PKCS7",
			plaintext:   "Hello, World! This is a test string for AES encryption.",
			key:         DefaultAESKey,
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherModeCBC,
			expectError: false,
		},
		{
			name:        "CustomKey_CBC_PKCS7",
			plaintext:   "Another test with a custom key.",
			key:         []byte("thisisasecretkey"), // 16 bytes
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherModeCBC,
			expectError: false,
		},
		{
			name:        "PKCS5_Padding",
			plaintext:   "Short string",
			key:         DefaultAESKey,
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS5,
			cipherMode:  CipherModeCBC,
			expectError: false,
		},
		{
			name:        "Invalid_Key_Length",
			plaintext:   "Test with invalid key",
			key:         []byte("short"), // Invalid key length
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherModeCBC,
			expectError: true,
		},
		{
			name:        "Invalid_IV_Length",
			plaintext:   "Test with invalid IV",
			key:         DefaultAESKey,
			iv:          []byte("shortiv"), // Invalid IV length
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherModeCBC,
			expectError: true,
		},
		{
			name:        "Unsupported_Padding_Mode",
			plaintext:   "Test unsupported padding",
			key:         DefaultAESKey,
			iv:          DefaultAESIV,
			paddingMode: PaddingMode(99), // Unsupported padding
			cipherMode:  CipherModeCBC,
			expectError: true,
		},
		{
			name:        "Unsupported_Cipher_Mode",
			plaintext:   "Test unsupported cipher",
			key:         DefaultAESKey,
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherMode(99), // Unsupported cipher
			expectError: true,
		},
		{
			name:        "ECB_Not_Implemented",
			plaintext:   "Test ECB",
			key:         DefaultAESKey,
			iv:          DefaultAESIV,
			paddingMode: PaddingModePKCS7,
			cipherMode:  CipherModeECB,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aes := NewAES().WithKey(tt.key).WithIV(tt.iv).WithMode(tt.paddingMode, tt.cipherMode)

			// Test EncryptToBase64 and DecryptFromBase64
			encryptedBase64, err := aes.EncryptToBase64(tt.plaintext)
			if (err != nil) != tt.expectError {
				t.Fatalf("EncryptToBase64() error = %v, expectError %v", err, tt.expectError)
			}
			if tt.expectError {
				return
			}

			decryptedBase64, err := aes.DecryptFromBase64(encryptedBase64)
			if err != nil {
				t.Fatalf("DecryptFromBase64() error = %v", err)
			}
			if decryptedBase64 != tt.plaintext {
				t.Errorf("Base64 encryption/decryption mismatch\nGot:  %s\nWant: %s", decryptedBase64, tt.plaintext)
			}

			// Test EncryptPlainText and DecryptCipherText
			encryptedHex, err := aes.EncryptPlainText(tt.plaintext)
			if (err != nil) != tt.expectError {
				t.Fatalf("EncryptPlainText() error = %v, expectError %v", err, tt.expectError)
			}
			if tt.expectError {
				return
			}

			decryptedHex, err := aes.DecryptCipherText(encryptedHex)
			if err != nil {
				t.Fatalf("DecryptCipherText() error = %v", err)
			}
			if decryptedHex != tt.plaintext {
				t.Errorf("Hex encryption/decryption mismatch\nGot:  %s\nWant: %s", decryptedHex, tt.plaintext)
			}
		})
	}
}

func TestAESCrypto_Encrypt_ErrorCases(t *testing.T) {
	aes := NewAES()
	// Test invalid IV length
	_, err := aes.WithIV([]byte{1, 2, 3}).Encrypt([]byte("test"))
	if err == nil || err.Error() != "IV length (3) must be equal to block size (16)" {
		t.Errorf("Expected IV length error, got %v", err)
	}

	// Test unsupported padding mode
	aes = NewAES().WithMode(PaddingMode(99), CipherModeCBC)
	_, err = aes.Encrypt([]byte("test"))
	if err == nil || err.Error() != "unsupported padding mode: 99" {
		t.Errorf("Expected unsupported padding mode error, got %v", err)
	}

	// Test unsupported cipher mode
	aes = NewAES().WithMode(PaddingModePKCS7, CipherMode(99))
	_, err = aes.Encrypt([]byte("test"))
	if err == nil || err.Error() != "unsupported cipher mode: 99" {
		t.Errorf("Expected unsupported cipher mode error, got %v", err)
	}

	// Test ECB not implemented
	aes = NewAES().WithMode(PaddingModePKCS7, CipherModeECB)
	_, err = aes.Encrypt([]byte("test"))
	if err == nil || err.Error() != "ECB cipher mode not implemented" {
		t.Errorf("Expected ECB not implemented error, got %v", err)
	}
}

func TestAESCrypto_Decrypt_ErrorCases(t *testing.T) {
	aes := NewAES()
	// Test invalid IV length
	_, err := aes.WithIV([]byte{1, 2, 3}).Decrypt([]byte("testtesttesttest")) // Must be multiple of block size
	if err == nil || err.Error() != "IV length (3) must be equal to block size (16)" {
		t.Errorf("Expected IV length error, got %v", err)
	}

	// Test ciphertext not multiple of block size
	aes = NewAES()
	_, err = aes.Decrypt([]byte("short"))
	if err == nil || err.Error() != "ciphertext is not a multiple of the block size" {
		t.Errorf("Expected ciphertext length error, got %v", err)
	}

	// Test unsupported padding mode
	aes = NewAES().WithMode(PaddingMode(99), CipherModeCBC)
	_, err = aes.Decrypt([]byte("testtesttesttest"))
	if err == nil || err.Error() != "unsupported padding mode: 99" {
		t.Errorf("Expected unsupported padding mode error, got %v", err)
	}

	// Test unsupported cipher mode
	aes = NewAES().WithMode(PaddingModePKCS7, CipherMode(99))
	_, err = aes.Decrypt([]byte("testtesttesttest"))
	if err == nil || err.Error() != "unsupported cipher mode: 99" {
		t.Errorf("Expected unsupported cipher mode error, got %v", err)
	}

	// Test ECB not implemented
	aes = NewAES().WithMode(PaddingModePKCS7, CipherModeECB)
	_, err = aes.Decrypt([]byte("testtesttesttest"))
	if err == nil || err.Error() != "ECB cipher mode not implemented" {
		t.Errorf("Expected ECB not implemented error, got %v", err)
	}

	// Test PKCS7 unpadding error (e.g., corrupted data)
	aes = NewAES().WithMode(PaddingModePKCS7, CipherModeCBC)
	// Create a valid ciphertext, then corrupt it slightly to trigger unpadding error
	validPlaintext := []byte("some data to encrypt")
	encrypted, _ := aes.Encrypt(validPlaintext)
	// Corrupt the last byte to make unpadding fail
	encrypted[len(encrypted)-1] = 0xFF
	_, err = aes.Decrypt(encrypted)
	if err == nil || err.Error() != "PKCS7 unpadding failed: invalid padding on input" {
		t.Errorf("Expected PKCS7 unpadding error, got %v", err)
	}
}

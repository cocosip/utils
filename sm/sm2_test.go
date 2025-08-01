package sm

import (
	"crypto/rand"
	"testing"

	"github.com/tjfoc/gmsm/sm2"
)

// TestSM2Crypto covers the core functionalities: key generation, encryption/decryption, and signing/verification.
func TestSM2Crypto(t *testing.T) {
	// 1. Test Key Generation
	privateKeyHex, publicKeyHex, err := NewSM2Key(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate SM2 key: %v", err)
	}
	if privateKeyHex == "" || publicKeyHex == "" {
		t.Fatal("Generated private or public key should not be empty")
	}

	// 2. Test Encryption and Decryption
	plaintext := []byte("hello world, this is a test message")

	// Test with C1C2C3 mode
	t.Run("EncryptDecrypt_C1C2C3", func(t *testing.T) {
		// Encrypt
		ciphertextHex, err := EncryptSM2(publicKeyHex, plaintext, rand.Reader, sm2.C1C2C3)
		if err != nil {
			t.Fatalf("Encryption failed with C1C2C3 mode: %v", err)
		}

		// Decrypt
		decryptedText, err := DecryptSM2(privateKeyHex, ciphertextHex, sm2.C1C2C3)
		if err != nil {
			t.Fatalf("Decryption failed with C1C2C3 mode: %v", err)
		}

		// Verify
		if string(plaintext) != decryptedText {
			t.Errorf("Decrypted text does not match original plaintext. got: %s, want: %s", decryptedText, string(plaintext))
		}
	})

	// Test with C1C3C2 mode
	t.Run("EncryptDecrypt_C1C3C2", func(t *testing.T) {
		// Encrypt
		ciphertextHex, err := EncryptSM2(publicKeyHex, plaintext, rand.Reader, sm2.C1C3C2)
		if err != nil {
			t.Fatalf("Encryption failed with C1C3C2 mode: %v", err)
		}

		// Decrypt
		decryptedText, err := DecryptSM2(privateKeyHex, ciphertextHex, sm2.C1C3C2)
		if err != nil {
			t.Fatalf("Decryption failed with C1C3C2 mode: %v", err)
		}

		// Verify
		if string(plaintext) != decryptedText {
			t.Errorf("Decrypted text does not match original plaintext. got: %s, want: %s", decryptedText, string(plaintext))
		}
	})

	// 3. Test Signing and Verification
	msg := []byte("this is a test message for signing")

	t.Run("SignVerify", func(t *testing.T) {
		// Sign
		signatureHex, err := SignSM2(privateKeyHex, msg, rand.Reader)
		if err != nil {
			t.Fatalf("Signing failed: %v", err)
		}

		// Verify a valid signature
		isValid, err := VerifySM2(publicKeyHex, msg, signatureHex)
		if err != nil {
			t.Fatalf("Verification of a valid signature failed: %v", err)
		}
		if !isValid {
			t.Error("Signature should be valid, but verification returned false")
		}

		// Verify an invalid signature (using a different message)
		wrongMsg := []byte("this is a different message")
		isValid, err = VerifySM2(publicKeyHex, wrongMsg, signatureHex)
		if err != nil {
			t.Fatalf("Verification with a wrong message failed: %v", err)
		}
		if isValid {
			t.Error("Signature should be invalid for a different message, but verification returned true")
		}
	})
}

// TestErrorCases checks for expected errors with invalid inputs.
func TestErrorCases(t *testing.T) {
	// Generate a valid key pair for testing
	privateKeyHex, publicKeyHex, err := NewSM2Key(rand.Reader)
	if err != nil {
		t.Fatalf("Key generation failed during setup for error case tests: %v", err)
	}

	t.Run("UnsupportedEncryptMode", func(t *testing.T) {
		plaintext := []byte("test")
		unsupportedMode := 99 // An invalid mode
		_, err := EncryptSM2(publicKeyHex, plaintext, rand.Reader, unsupportedMode)
		if err == nil {
			t.Error("Expected an error for unsupported encryption mode, but got nil")
		}
	})

	t.Run("UnsupportedDecryptMode", func(t *testing.T) {
		ciphertextHex := "dummyciphertext"
		unsupportedMode := 99 // An invalid mode
		_, err := DecryptSM2(privateKeyHex, ciphertextHex, unsupportedMode)
		if err == nil {
			t.Error("Expected an error for unsupported decryption mode, but got nil")
		}
	})

	t.Run("InvalidPublicKey", func(t *testing.T) {
		_, err := VerifySM2("invalid-hex-key", []byte("msg"), "invalid-signature")
		if err == nil {
			t.Error("Expected an error for an invalid public key, but got nil")
		}
	})

	t.Run("InvalidPrivateKey", func(t *testing.T) {
		_, err := DecryptSM2("invalid-hex-key", "dummyciphertext", sm2.C1C2C3)
		if err == nil {
			t.Error("Expected an error for an invalid private key, but got nil")
		}
	})
}

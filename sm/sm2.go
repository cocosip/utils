// Package sm provides SM2, SM3, and SM4 cryptographic functionalities.
// It leverages the 'github.com/tjfoc/gmsm' library.
package sm

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	gmsm_sm2 "github.com/tjfoc/gmsm/sm2"
	gmsm_x509 "github.com/tjfoc/gmsm/x509"
)

// NewSM2Key generates a new SM2 private and public key pair.
// It returns the private key and public key as hex-encoded strings.
// The `random` parameter is an `io.Reader` for cryptographic randomness; if nil, `crypto/rand.Reader` is used.
func NewSM2Key(random io.Reader) (privateKeyHex string, publicKeyHex string, err error) {
	if random == nil {
		random = rand.Reader
	}
	priv, err := gmsm_sm2.GenerateKey(random)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate SM2 key: %w", err)
	}
	// Convert private and public keys to hex-encoded strings.
	return gmsm_x509.WritePrivateKeyToHex(priv), gmsm_x509.WritePublicKeyToHex(&priv.PublicKey), nil
}

// EncryptSM2 encrypts plaintext using the SM2 public key.
// `pub` is the hex-encoded SM2 public key.
// `plaintext` is the data to be encrypted (raw bytes).
// `random` is an `io.Reader` for cryptographic randomness; if nil, `crypto/rand.Reader` is used.
// `mode` specifies the encryption mode: sm2.C1C2C3 or sm2.C1C3C2.
// It returns the hex-encoded ciphertext.
func EncryptSM2(pub string, plaintext []byte, random io.Reader, mode int) (ciphertextHex string, err error) {
	pubK, err := gmsm_x509.ReadPublicKeyFromHex(pub)
	if err != nil {
		return "", fmt.Errorf("failed to read public key from hex: %w", err)
	}

	if random == nil {
		random = rand.Reader
	}

	if mode != gmsm_sm2.C1C2C3 && mode != gmsm_sm2.C1C3C2 {
		return "", fmt.Errorf("unsupported SM2 encryption mode: %d", mode)
	}

	// Encrypt the plaintext using the specified mode.
	cipherBuf, err := gmsm_sm2.Encrypt(pubK, plaintext, random, mode) // plaintext is now []byte
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}
	return hex.EncodeToString(cipherBuf), nil
}

// DecryptSM2 decrypts ciphertext using the SM2 private key.
// `priv` is the hex-encoded SM2 private key.
// `ciphertextHex` is the hex-encoded encrypted data.
// `mode` specifies the decryption mode: sm2.C1C2C3 or sm2.C1C3C2.
// It returns the original plaintext as a byte slice.
func DecryptSM2(priv string, ciphertextHex string, mode int) (plaintext []byte, err error) {
	privK, err := gmsm_x509.ReadPrivateKeyFromHex(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key from hex: %w", err)
	}

	ciphertextBytes, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext from hex: %w", err)
	}

	if mode != gmsm_sm2.C1C2C3 && mode != gmsm_sm2.C1C3C2 {
		return nil, fmt.Errorf("unsupported SM2 decryption mode: %d", mode)
	}

	// Decrypt the ciphertext using the specified mode.
	plainBuf, err := gmsm_sm2.Decrypt(privK, ciphertextBytes, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return plainBuf, nil
}

// SignSM2 signs a message using the SM2 private key.
// `priv` is the hex-encoded SM2 private key.
// `msg` is the message to be signed (raw bytes).
// `random` is an `io.Reader` for cryptographic randomness; if nil, `crypto/rand.Reader` is used.
// It returns the hex-encoded signature.
func SignSM2(priv string, msg []byte, random io.Reader) (signatureHex string, err error) {
	privK, err := gmsm_x509.ReadPrivateKeyFromHex(priv)
	if err != nil {
		return "", fmt.Errorf("failed to read private key from hex for signing: %w", err)
	}

	if random == nil {
		random = rand.Reader
	}

	// Sign the message.
	sig, err := privK.Sign(random, msg, nil) // msg is already []byte
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %w", err)
	}
	return hex.EncodeToString(sig), nil
}

// VerifySM2 verifies an SM2 signature.
// `pub` is the hex-encoded SM2 public key.
// `msg` is the original message that was signed (raw bytes).
// `signatureHex` is the hex-encoded signature data.
// It returns true if the signature is valid, false otherwise, and an error if any issue occurs during verification.
func VerifySM2(pub string, msg []byte, signatureHex string) (isValid bool, err error) {
	pubK, err := gmsm_x509.ReadPublicKeyFromHex(pub)
	if err != nil {
		return false, fmt.Errorf("failed to read public key from hex for verification: %w", err)
	}

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature from hex: %w", err)
	}

	// Verify the signature.
	isValid = pubK.Verify(msg, signature)
	return isValid, nil
}
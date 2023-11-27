package sm

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"io"
)

func NewSM2Key(random io.Reader) (string, string, error) {
	if random == nil {
		random = rand.Reader
	}
	priv, err := sm2.GenerateKey(random)
	if err != nil {
		return "", "", err
	}
	return x509.WritePrivateKeyToHex(priv), x509.WritePublicKeyToHex(&priv.PublicKey), nil
}

func EncryptSM2(pub string, plaintext string, random io.Reader, mode int) (string, error) {
	pubK, err := x509.ReadPublicKeyFromHex(pub)
	if err != nil {
		return "", err
	}

	if random == nil {
		random = rand.Reader
	}

	//sm2.C1C2C3
	cipherBuf, err := sm2.Encrypt(pubK, []byte(plaintext), random, mode)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipherBuf), nil
}

func DecryptSM2(priv string, ciphertext string, mode int) (string, error) {
	privK, err := x509.ReadPrivateKeyFromHex(priv)
	if err != nil {
		return "", err
	}

	b, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	//sm2.C1C2C3
	plainBuf, err := sm2.Decrypt(privK, b, mode)
	if err != nil {
		return "", err
	}
	return string(plainBuf), nil
}

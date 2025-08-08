package encoding

import "encoding/base64"

// EncodeBase64 encodes a string to base64.
// param: data string to encode
// return: base64 encoded string
func EncodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// EncodeBase64Bytes encodes a byte slice to base64.
// param: data byte slice to encode
// return: base64 encoded string
func EncodeBase64Bytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 string to string.
// param: data base64 encoded string
// return: decoded string, error if decode fails
func DecodeBase64(data string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DecodeBase64ToBytes decodes a base64 string to byte slice.
// param: data base64 encoded string
// return: decoded byte slice, error if decode fails
func DecodeBase64ToBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

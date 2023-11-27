package encoding

import "encoding/base64"

func EncodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeBase64(data string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

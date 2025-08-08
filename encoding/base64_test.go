package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	v1 := EncodeBase64("Hello World")
	assert.Equal(t, "SGVsbG8gV29ybGQ=", v1, "encode base64 fail")
}

func TestEncodeBase64Bytes(t *testing.T) {
	v := EncodeBase64Bytes([]byte("GoLang测试"))
	assert.Equal(t, "R29MYW5n5rWL6K+V", v, "encode base64 bytes fail")
}

func TestDecodeBase64(t *testing.T) {
	v1, err := DecodeBase64("6L+Z5piv5Lit5Zu9")
	assert.NoErrorf(t, err, "decode base64 error %s", err)
	assert.Equal(t, "这是中国", v1, "decode base64 fail")
}

func TestDecodeBase64ToBytes(t *testing.T) {
	v, err := DecodeBase64ToBytes("SGVsbG8gV29ybGQ=")
	assert.NoError(t, err, "decode base64 to bytes error")
	assert.Equal(t, []byte("Hello World"), v, "decode base64 to bytes fail")
}

func TestDecodeBase64_Error(t *testing.T) {
	_, err := DecodeBase64("!!!notbase64!!!")
	assert.Error(t, err, "should error on invalid base64 string")
}

func TestDecodeBase64ToBytes_Error(t *testing.T) {
	_, err := DecodeBase64ToBytes("!!!notbase64!!!")
	assert.Error(t, err, "should error on invalid base64 string")
}

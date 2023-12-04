package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	v1 := EncodeBase64("Hello World")
	assert.Equal(t, "SGVsbG8gV29ybGQ=", v1, "encode base64 fail")
}

func TestDecodeBase64(t *testing.T) {
	v1, err := DecodeBase64("6L+Z5piv5Lit5Zu9")
	assert.NoErrorf(t, err, "decode base64 error %s", err)
	assert.Equal(t, "这是中国", v1, "decode base64 fail")
}

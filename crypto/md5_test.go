package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMD5(t *testing.T) {
	v1 := MD5("HelloWorld")
	assert.Equal(t, "68e109f0f40ca72a15e05cc22786f8e6", v1)
}

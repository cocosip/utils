package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSHA1(t *testing.T) {
	v1 := SHA1("This is a test")
	assert.Equal(t, "a54d88e06612d820bc3be72877c74f257b561b19", v1)
}

func TestSHA256(t *testing.T) {
	v1 := SHA256("This is a test")
	assert.Equal(t, "c7be1ed902fb8dd4d48997c6452f5d7e509fbcdbe2808b16bcf4edce4c07d14e", v1)
}

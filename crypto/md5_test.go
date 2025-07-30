package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMD5(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "Simple string",
			input:    "HelloWorld",
			expected: "68e109f0f40ca72a15e05cc22786f8e6",
		},
		{
			name:     "String with spaces",
			input:    "Hello World",
			expected: "b10a8db164e0754105b7a99be72e3fe5",
		},
		{
			name:     "String with special characters",
			input:    "!@#$%^&*()_+",
			expected: "04dde9f462255fe14b5160bbf2acffe8",
		},
		{
			name:     "Longer string",
			input:    "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.",
			expected: "f168d89e05b664041ee6745f050caa4b",
		},
		{
			name:     "UTF-8 characters",
			input:    "你好世界",
			expected: "65396ee4aad0b4f17aacd1c6112ee364",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MD5(tt.input)
			assert.Equal(t, tt.expected, actual, "MD5 hash mismatch for input: %s", tt.input)
		})
	}
}

package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSHA1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			name:     "Simple string",
			input:    "HelloWorld",
			expected: "db8ac1c259eb89d4a131b253bacfca5f319d54f2",
		},
		{
			name:     "String with spaces",
			input:    "Hello World",
			expected: "0a4d55a8d778e5022fab701977c5d840bbc486d0",
		},
		{
			name:     "String with special characters",
			input:    "!@#$%^&*()_+",
			expected: "d0b9abafaf5a393954f53e47715c833f0c18075d",
		},
		{
			name:     "Longer string",
			input:    "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.",
			expected: "8eb09e076e6afa8aaaf8ff172f76cee84c791994",
		},
		{
			name:     "UTF-8 characters",
			input:    "你好世界",
			expected: "dabaa5fe7c47fb21be902480a13013f16a1ab6eb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SHA1(tt.input)
			assert.Equal(t, tt.expected, actual, "SHA1 hash mismatch for input: %s", tt.input)
		})
	}
}

func TestSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "Simple string",
			input:    "HelloWorld",
			expected: "872e4e50ce9990d8b041330c47c9ddd11bec6b503ae9386a99da8584e9bb12c4",
		},
		{
			name:     "String with spaces",
			input:    "Hello World",
			expected: "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e",
		},
		{
			name:     "String with special characters",
			input:    "!@#$%^&*()_+",
			expected: "36d3e1bc65f8b67935ae60f542abef3e55c5bbbd547854966400cc4f022566cb",
		},
		{
			name:     "Longer string",
			input:    "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.",
			expected: "635241ac823ee4a81fbb410c92be616b0a89191083d8d7b5d232c823dc8df4f5",
		},
		{
			name:     "UTF-8 characters",
			input:    "你好世界",
			expected: "beca6335b20ff57ccc47403ef4d9e0b8fccb4442b3151c2e7d50050673d43172",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SHA256(tt.input)
			assert.Equal(t, tt.expected, actual, "SHA256 hash mismatch for input: %s", tt.input)
		})
	}
}

func TestSHA512(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			name:     "Simple string",
			input:    "HelloWorld",
			expected: "8ae6ae71a75d3fb2e0225deeb004faf95d816a0a58093eb4cb5a3aa0f197050d7a4dc0a2d5c6fbae5fb5b0d536a0a9e6b686369fa57a027687c3630321547596",
		},
		{
			name:     "String with spaces",
			input:    "Hello World",
			expected: "2c74fd17edafd80e8447b0d46741ee243b7eb74dd2149a0ab1b9246fb30382f27e853d8585719e0e67cbda0daa8f51671064615d645ae27acb15bfb1447f459b",
		},
		{
			name:     "String with special characters",
			input:    "!@#$%^&*()_+",
			expected: "ef38ee69ccb0ae69dfd99fdd896a63067c0393c3ea41548746c2ec103a3d77538367117600ff0e550a3ad67787b6072f2fb989487474b55d62c6508d809e44e1",
		},
		{
			name:     "Longer string",
			input:    "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.",
			expected: "eadacad7da6a99f7058fcba27e5ee6986de052ed4e2a24475b88acb863fe95a6bfe996e1523d86c085dbe45655f96d9f85ad10e2cb6e93d22500c59bc2ae4a34",
		},
		{
			name:     "UTF-8 characters",
			input:    "你好世界",
			expected: "4b28a152c8e203ebb52e099301041e3cf704a56190d3097ec8b086a0f9bfb4b9d533ce71fc3bcf374359e506dc5f17322ec3911eac8dd8f5b35308d938ba0c26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SHA512(tt.input)
			assert.Equal(t, tt.expected, actual, "SHA512 hash mismatch for input: %s", tt.input)
		})
	}
}

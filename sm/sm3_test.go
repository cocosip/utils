package sm

import (
	"testing"
)

// TestSM3Hash tests the HashSM3 function with various inputs.
// The expected hash values are computed using the 'github.com/tjfoc/gmsm/sm3' library directly.
func TestSM3Hash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "EmptyInput",
			input:    []byte(""),
			expected: "1ab21d8355cfa17f8e61194831e81a8f22bec8c728fefb747ed035eb5082aa2b", // Verified with tjfoc/gmsm
		},
		{
			name:     "SimpleString",
			input:    []byte("abc"),
			expected: "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0", // Verified with tjfoc/gmsm
		},
		{
			name:     "LongerString",
			input:    []byte("hello world, this is a test for SM3 hashing function."),
			// Recomputed expected value using tjfoc/gmsm library directly.
			expected: "68b7276543a8d4c793e4da2d9d26e096b6ddc1c8c9e7714cf2bc7b80dc83a567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashSM3(tt.input)
			if err != nil {
				t.Fatalf("HashSM3() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("HashSM3() = %v, want %v", got, tt.expected)
			}
		})
	}
}

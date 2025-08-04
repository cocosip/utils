package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntOrDefault(t *testing.T) {
	assert.Equal(t, 123, ParseIntOrDefault("123"))
	assert.Equal(t, -45, ParseIntOrDefault("-45"))
	assert.Equal(t, 0, ParseIntOrDefault("abc"))
	assert.Equal(t, 0, ParseIntOrDefault("12.3"))
}

func TestParseIntE(t *testing.T) {
	v, err := ParseIntE("123")
	assert.NoError(t, err)
	assert.Equal(t, 123, v)

	_, err = ParseIntE("abc")
	assert.Error(t, err)
}

func TestParseUintOrDefault(t *testing.T) {
	assert.Equal(t, uint(123), ParseUintOrDefault("123"))
	assert.Equal(t, uint(0), ParseUintOrDefault("-45")) // Negative is an error
	assert.Equal(t, uint(0), ParseUintOrDefault("abc"))
}

func TestParseUintE(t *testing.T) {
	v, err := ParseUintE("123")
	assert.NoError(t, err)
	assert.Equal(t, uint(123), v)

	_, err = ParseUintE("-45")
	assert.Error(t, err, "Should error on negative numbers")
}

func TestParseInt64OrDefault(t *testing.T) {
	assert.Equal(t, int64(1234567890), ParseInt64OrDefault("1234567890"))
	assert.Equal(t, int64(-1234567890), ParseInt64OrDefault("-1234567890"))
	assert.Equal(t, int64(0), ParseInt64OrDefault("abc"))
}

func TestParseInt64E(t *testing.T) {
	v, err := ParseInt64E("1234567890")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234567890), v)

	_, err = ParseInt64E("abc")
	assert.Error(t, err)
}

func TestParseUint64OrDefault(t *testing.T) {
	assert.Equal(t, uint64(123456789012345), ParseUint64OrDefault("123456789012345"))
	assert.Equal(t, uint64(0), ParseUint64OrDefault("-123"))
	assert.Equal(t, uint64(0), ParseUint64OrDefault("abc"))
}

func TestParseUint64E(t *testing.T) {
	v, err := ParseUint64E("123456789012345")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456789012345), v)

	_, err = ParseUint64E("-123")
	assert.Error(t, err)
}

func TestAnonymizeString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		max     int
		rep     rune
		expected string
	}{
		{"long string", "13812345678", 5, '*', "1***8"},
		{"short string", "abc", 5, '#', "a###c"},
		{"two char string", "ab", 4, '*', "a**b"},
		{"one char string", "a", 4, '*', "a***"},
		{"empty string", "", 5, '*', ""},
		{
		name:     "max less than 2",
		input:    "abcdef",
		max:      1,
		rep:      '*',
		expected: "af", // max is corrected to 2, showing first and last char
	},
		{"unicode string", "⽜⼈是很棒的", 6, '-', "⽜----的"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, AnonymizeString(tt.input, tt.max, tt.rep))
		})
	}
}
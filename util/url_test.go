package util

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid HTTPS", "https://example.com", true},
		{"Valid HTTP", "http://example.com/path?query=1", true},
		{"Valid with port", "ftp://user:pass@host:21/path", true},
		{"No scheme", "example.com", false},
		{"Relative path", "/a/b/c", false},
		{"Just a string", "not-a-url", false},
		{"Empty string", "", false},
		{"Scheme only", "http://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsURL(tt.input))
		})
	}
}

func TestAddURLParams(t *testing.T) {
	t.Run("Add to URL without params", func(t *testing.T) {
		rawURL := "http://example.com/path"
		params := map[string]string{"a": "1", "b": "2"}
		newURL, err := AddURLParams(rawURL, params)
		assert.NoError(t, err)
		// The order of params is not guaranteed, so we check for both possibilities
		assert.Contains(t, []string{"http://example.com/path?a=1&b=2", "http://example.com/path?b=2&a=1"}, newURL)
	})

	t.Run("Add to URL with existing params", func(t *testing.T) {
		rawURL := "http://example.com?a=1"
		params := map[string]string{"b": "2", "c": "3"}
		newURL, err := AddURLParams(rawURL, params)
		assert.NoError(t, err)
		// We can't guarantee order, so we parse it back and check values
		u, _ := url.Parse(newURL)
		q := u.Query()
		assert.Equal(t, "1", q.Get("a"))
		assert.Equal(t, "2", q.Get("b"))
		assert.Equal(t, "3", q.Get("c"))
	})

	t.Run("Overwrite existing param", func(t *testing.T) {
		rawURL := "http://example.com?a=1"
		params := map[string]string{"a": "new_value"}
		newURL, err := AddURLParams(rawURL, params)
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com?a=new_value", newURL)
	})

	t.Run("Invalid URL input", func(t *testing.T) {
		rawURL := "://invalid"
		params := map[string]string{"a": "1"}
		_, err := AddURLParams(rawURL, params)
		assert.Error(t, err)
	})
}

func TestGetURLParam(t *testing.T) {
	rawURL := "http://example.com?a=1&b=hello&c=&d"

	t.Run("Get existing param", func(t *testing.T) {
		val, ok := GetURLParam(rawURL, "a")
		assert.True(t, ok)
		assert.Equal(t, "1", val)
	})

	t.Run("Get another existing param", func(t *testing.T) {
		val, ok := GetURLParam(rawURL, "b")
		assert.True(t, ok)
		assert.Equal(t, "hello", val)
	})

	t.Run("Get param with empty value", func(t *testing.T) {
		val, ok := GetURLParam(rawURL, "c")
		assert.True(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("Get param with no value", func(t *testing.T) {
		val, ok := GetURLParam(rawURL, "d")
		assert.True(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("Get non-existent param", func(t *testing.T) {
		val, ok := GetURLParam(rawURL, "z")
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("Invalid URL input", func(t *testing.T) {
		val, ok := GetURLParam("://invalid", "a")
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})
}

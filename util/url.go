// Package util provides a collection of utility functions.
package util

import (
	"net/url"
)

// IsURL checks if a string is a valid, absolute URL.
// It requires the URL to have a scheme (e.g., "http", "https") and a host.
func IsURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// AddURLParams adds query parameters to a raw URL string.
// It correctly handles URLs that already have query parameters.
func AddURLParams(rawURL string, params map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// GetURLParam retrieves a single query parameter value from a raw URL string.
// It returns the value and a boolean indicating whether the parameter was found.
func GetURLParam(rawURL string, key string) (string, bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}
	val := u.Query().Get(key)
	// u.Query().Get returns "" for non-existent keys, so we check with Has
	return val, u.Query().Has(key)
}
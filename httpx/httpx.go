// Package httpx provides utility functions for working with HTTP in Go.
// It includes helpers for creating custom HTTP clients, handling responses,
// and managing headers.
package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultHttpClientTimeout defines the default timeout for the HTTP client.
	// This value is used when no custom timeout is specified.
	DefaultHttpClientTimeout = time.Second * 100
)

// TransportOption is a function type used to configure an http.Transport.
// It follows the functional options pattern, allowing for flexible transport customization.
type TransportOption func(transport *http.Transport) error

// ClientOption is a function type used to configure an http.Client.
// It follows the functional options pattern, enabling modular and readable client configuration.
type ClientOption func(*http.Client) error

// WithTimeout returns a ClientOption that sets the timeout for the http.Client.
// The duration specifies the total time limit for a request, including connection time,
// reading headers, and reading the response body.
func WithTimeout(duration time.Duration) ClientOption {
	return func(client *http.Client) error {
		client.Timeout = duration
		return nil
	}
}

// WithTransport returns a ClientOption that configures the http.Client's transport.
// It takes one or more TransportOption functions to customize the transport settings.
// This function clones the http.DefaultTransport to avoid modifying the global default.
func WithTransport(opts ...TransportOption) ClientOption {
	return func(client *http.Client) error {
		// Clone the default transport to avoid side effects on the global http.DefaultTransport.
		transport := http.DefaultTransport.(*http.Transport).Clone()
		for _, opt := range opts {
			if err := opt(transport); err != nil {
				return fmt.Errorf("failed to apply transport option: %w", err)
			}
		}
		client.Transport = transport
		return nil
	}
}

// WithCheckRedirect returns a ClientOption that sets the redirect policy for the http.Client.
// The provided function determines how to handle redirects. If the function returns an error,
// the redirect is not followed.
func WithCheckRedirect(f func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *http.Client) error {
		client.CheckRedirect = f
		return nil
	}
}

// WithCookieJar returns a ClientOption that sets the cookie jar for the http.Client.
// The cookie jar is used to store and send cookies with requests.
func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(client *http.Client) error {
		client.Jar = jar
		return nil
	}
}

// WithTransportProxyURL returns a TransportOption that configures the HTTP proxy for the transport.
// It parses the proxyURL and sets it on the http.Transport. If a username and password are provided,
// they are included for proxy authentication.
func WithTransportProxyURL(proxyURL string, username string, password string) TransportOption {
	return func(transport *http.Transport) error {
		if proxyURL == "" {
			return nil // No proxy to set
		}

		u, err := url.Parse(proxyURL)
		if err != nil {
			return fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		if strings.TrimSpace(username) != "" && strings.TrimSpace(password) != "" {
			u.User = url.UserPassword(username, password)
		}

		transport.Proxy = http.ProxyURL(u)
		return nil
	}
}

// NewHttpClient creates a new http.Client with the specified options.
// It starts with a default timeout and applies any provided ClientOption functions.
// This allows for easy and flexible construction of custom HTTP clients.
func NewHttpClient(opts ...ClientOption) (*http.Client, error) {
	client := &http.Client{
		Timeout: DefaultHttpClientTimeout,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	return client, nil
}

// SetBearerAuth sets the Authorization header on an http.Request to use Bearer authentication.
// It formats the token and adds it to the request header.
func SetBearerAuth(request *http.Request, token string) {
	bearer := fmt.Sprintf("Bearer %s", token)
	request.Header.Set("Authorization", bearer)
}

// GetForwardedFor parses the "X-Forwarded-For" header and returns a slice of IP addresses.
// This header is commonly used by proxies and load balancers to indicate the original client IP.
// The function handles comma-separated lists of IPs and trims whitespace.
func GetForwardedFor(header http.Header) []string {
	val := header.Get("X-Forwarded-For")
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	hosts := make([]string, 0, len(parts))
	for _, part := range parts {
		host := strings.TrimSpace(part)
		if host != "" {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

// ReadResponseAsJson reads the response body, ensures a success status code (2xx),
// and unmarshals the JSON content into the provided interface.
// It automatically closes the response body. If the status code is not successful
// or if JSON unmarshaling fails, it returns an error.
func ReadResponseAsJson(response *http.Response, resp interface{}) error {
	body, err := ReadResponseAsBytes(response)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, resp); err != nil {
		// Include the response body in the error for better debugging.
		return fmt.Errorf("json unmarshal failed: %w; body: %s", err, string(body))
	}
	return nil
}

// ReadResponseAsBytes reads the response body and ensures a success status code (2xx).
// It automatically closes the response body. If the status code is not in the 200-299 range,
// it returns an error that includes the response body for context.
func ReadResponseAsBytes(response *http.Response) ([]byte, error) {
	defer CloseResponse(response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(body))
	}

	return body, nil
}

// CloseResponse closes the response body and handles any errors that occur.
// If callbacks are provided, they are called with the error. Otherwise, the error
// is logged using the standard slog logger. This helper prevents resource leaks
// from unclosed response bodies.
func CloseResponse(response *http.Response, callbacks ...func(error)) {
	if err := response.Body.Close(); err != nil {
		if len(callbacks) > 0 {
			for _, callback := range callbacks {
				if callback != nil {
					callback(err)
				}
			}
		} else {
			slog.Error("failed to close http response body", "error", err, "url", response.Request.URL.String())
		}
	}
}

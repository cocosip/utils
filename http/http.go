package http

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
	DefaultHttpClientTimeout = time.Second * 100
)

type TransportOption func(transport *http.Transport) error

type ClientOption func(*http.Client) error

func WithTimeout(duration time.Duration) ClientOption {
	return func(client *http.Client) error {
		client.Timeout = duration
		return nil
	}
}

func WithTransport(opts ...TransportOption) ClientOption {
	return func(client *http.Client) error {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		for _, opt := range opts {
			if err := opt(transport); err != nil {
				return err
			}

		}
		client.Transport = transport
		return nil
	}
}

func WithCheckRedirect(f func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *http.Client) error {
		client.CheckRedirect = f
		return nil
	}
}

func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(client *http.Client) error {
		client.Jar = jar
		return nil
	}
}

func WithTransportProxyURL(proxyURL string, username string, password string) TransportOption {
	return func(transport *http.Transport) error {
		if proxyURL != "" {
			u, err := url.Parse(proxyURL)
			if err != nil {
				return err
			}
			if strings.TrimSpace(username) != "" && strings.TrimSpace(password) != "" {
				u.User = url.UserPassword(username, password)
			}
			transport.Proxy = http.ProxyURL(u)
		}
		return nil
	}
}

func NewHttpClient(opts ...ClientOption) (*http.Client, error) {
	client := &http.Client{
		Timeout: DefaultHttpClientTimeout,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func SetBearerAuth(request *http.Request, token string) {
	bearer := fmt.Sprintf("Bearer %s", token)
	request.Header.Set("Authorization", bearer)
}

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

// ReadResponseAsJson reads the response body, ensures a success status code, and unmarshals it into the provided interface.
// It automatically closes the response body.
func ReadResponseAsJson(response *http.Response, resp interface{}) error {
	body, err := ReadResponseAsBytes(response)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("json unmarshal failed: %w; body: %s", err, string(body))
	}
	return nil
}

// ReadResponseAsBytes reads the response body and ensures a success status code (2xx).
// It automatically closes the response body.
// If the status code is not in the 200-299 range, it returns an error with the response body.
func ReadResponseAsBytes(response *http.Response) ([]byte, error) {
	defer CloseResponse(response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(body))
	}

	return body, nil
}

// CloseResponse closes the response body and logs an error if the close operation fails.
// Callbacks can be provided to handle the error manually.
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

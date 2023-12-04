package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
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
	var hosts []string
	arr := strings.Split(val, ",")
	for i := range arr {
		if strings.TrimSpace(arr[i]) != "" {
			hosts = append(hosts, strings.TrimSpace(arr[i]))
		}
	}
	return hosts
}

func ReadJsonResponse(response *http.Response, resp interface{}) error {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("return status code = %d; body = %s", response.StatusCode, string(b))
	}

	if err = json.Unmarshal(b, resp); err != nil {
		return fmt.Errorf("unmarshal error = %s; body = %s", err, string(b))
	}
	return nil
}

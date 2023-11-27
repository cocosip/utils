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

func NewProxyHttpClient(proxyURL string, username string, password string, timeouts ...time.Duration) (*http.Client, error) {
	timeout := DefaultHttpClientTimeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	cli := &http.Client{
		Timeout: timeout,
	}

	if strings.TrimSpace(proxyURL) != "" {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse proxyURL error: %s; proxyURL = %s ", err.Error(), proxyURL)
		}

		if strings.TrimSpace(username) != "" && strings.TrimSpace(password) != "" {
			u.User = url.UserPassword(username, password)
		}

		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.Proxy = http.ProxyURL(u)
		cli.Transport = transport
	}

	return cli, nil
}

func NewTimeoutHttpClient(timeouts ...time.Duration) *http.Client {
	timeout := DefaultHttpClientTimeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	return &http.Client{
		Timeout: timeout,
	}
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

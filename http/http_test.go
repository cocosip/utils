package http

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	_, err := NewHttpClient(
		WithTimeout(time.Second*20),
		WithTransport(
			WithTransportProxyURL("http://127.0.0.1", "admin", "123456"),
			func(transport *http.Transport) error {
				transport.IdleConnTimeout = time.Second * 100
				return nil
			},
		),
		WithCheckRedirect(func(req *http.Request, via []*http.Request) error {
			return nil
		}))
	if err != nil {
		t.Errorf("test new http client error -> %s", err.Error())
	}

	t.Logf("test new http client success")

}

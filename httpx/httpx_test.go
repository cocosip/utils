package httpx

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewHttpClient tests the NewHttpClient function with various options.
func TestNewHttpClient(t *testing.T) {
	// Test with no options
	client, err := NewHttpClient()
	if err != nil {
		t.Fatalf("NewHttpClient with no options failed: %v", err)
	}
	if client.Timeout != DefaultHttpClientTimeout {
		t.Errorf("Expected default timeout %v, got %v", DefaultHttpClientTimeout, client.Timeout)
	}

	// Test with timeout option
	customTimeout := time.Second * 5
	client, err = NewHttpClient(WithTimeout(customTimeout))
	if err != nil {
		t.Fatalf("NewHttpClient with timeout option failed: %v", err)
	}
	if client.Timeout != customTimeout {
		t.Errorf("Expected custom timeout %v, got %v", customTimeout, client.Timeout)
	}

	// Test with transport options (proxy and idle conn timeout)
	proxyURL := "http://127.0.0.1:8080"
	client, err = NewHttpClient(
		WithTransport(
			WithTransportProxyURL(proxyURL, "admin", "123456"),
			func(transport *http.Transport) error {
				transport.IdleConnTimeout = time.Second * 100
				return nil
			},
		),
	)
	if err != nil {
		t.Fatalf("NewHttpClient with transport options failed: %v", err)
	}

	if transport, ok := client.Transport.(*http.Transport); ok {
		if transport.Proxy == nil {
			t.Error("Expected proxy to be set, got nil")
		} else {
			proxy, _ := transport.Proxy(&http.Request{})
			if proxy.Host != "127.0.0.1:8080" {
				t.Errorf("Expected proxy host %s, got %s", "127.0.0.1:8080", proxy.Host)
			}
		}
		if transport.IdleConnTimeout != time.Second*100 {
			t.Errorf("Expected IdleConnTimeout %v, got %v", time.Second*100, transport.IdleConnTimeout)
		}
	} else {
		t.Error("Expected client transport to be *http.Transport")
	}

	// Test with CheckRedirect option
	redirectCalled := false
	client, err = NewHttpClient(
		WithCheckRedirect(func(req *http.Request, via []*http.Request) error {
			redirectCalled = true
			return http.ErrUseLastResponse // Prevent actual redirect
		}),
	)
	if err != nil {
		t.Fatalf("NewHttpClient with CheckRedirect option failed: %v", err)
	}

	// Simulate a redirect to check if the function is called
	_ = &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": []string{"/new-location"}}}
	req, _ := http.NewRequest("GET", "/", nil)
	client.CheckRedirect(req, []*http.Request{})
	if !redirectCalled {
		t.Error("Expected CheckRedirect function to be called, but it wasn't")
	}

	// Test with CookieJar option
	jar := &mockCookieJar{}
	client, err = NewHttpClient(WithCookieJar(jar))
	if err != nil {
		t.Fatalf("NewHttpClient with CookieJar option failed: %v", err)
	}
	if client.Jar != jar {
		t.Error("Expected custom CookieJar, got different one")
	}

	t.Logf("TestNewHttpClient completed successfully")
}

// mockCookieJar implements http.CookieJar for testing purposes.
type mockCookieJar struct{}

func (j *mockCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {}
func (j *mockCookieJar) Cookies(u *url.URL) []*http.Cookie             { return nil }

// TestSetBearerAuth tests the SetBearerAuth function.
func TestSetBearerAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	token := "test_token_123"

	SetBearerAuth(req, token)

	expectedHeader := "Bearer test_token_123"
	if req.Header.Get("Authorization") != expectedHeader {
		t.Errorf("Expected Authorization header %q, got %q", expectedHeader, req.Header.Get("Authorization"))
	}
}

// TestGetForwardedFor tests the GetForwardedFor function.
func TestGetForwardedFor(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		expected []string
	}{
		{
			name:     "Empty header",
			header:   http.Header{},
			expected: nil,
		},
		{
			name:     "Single IP",
			header:   http.Header{"X-Forwarded-For": []string{"192.168.1.1"}},
			expected: []string{"192.168.1.1"},
		},
		{
			name:     "Multiple IPs",
			header:   http.Header{"X-Forwarded-For": []string{"192.168.1.1, 10.0.0.1, 172.16.0.1"}},
			expected: []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"},
		},
		{
			name:     "IPs with spaces",
			header:   http.Header{"X-Forwarded-For": []string{"  192.168.1.1  ,   10.0.0.1  "}},
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
		{
			name:     "Empty parts",
			header:   http.Header{"X-Forwarded-For": []string{", ,192.168.1.1,"}},
			expected: []string{"192.168.1.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetForwardedFor(tt.header)
			if !compareStringSlices(got, tt.expected) {
				t.Errorf("GetForwardedFor() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

// compareStringSlices is a helper to compare two string slices.
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestReadResponseAsBytes tests the ReadResponseAsBytes function.
func TestReadResponseAsBytes(t *testing.T) {
	// Success case: 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test body content"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get response from test server: %v", err)
	}

	body, err := ReadResponseAsBytes(resp)
	if err != nil {
		t.Errorf("ReadResponseAsBytes failed for success case: %v", err)
	}
	if string(body) != "test body content" {
		t.Errorf("Expected body %q, got %q", "test body content", string(body))
	}

	// Failure case: 400 Bad Request
	serverErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request error"))
	}))
	defer serverErr.Close()

	respErr, err := http.Get(serverErr.URL)
	if err != nil {
		t.Fatalf("Failed to get error response from test server: %v", err)
	}

	_, err = ReadResponseAsBytes(respErr)
	if err == nil {
		t.Error("Expected an error for 400 status code, got nil")
	}
	expectedErr := "unexpected status code 400: bad request error"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

// TestReadResponseAsJson tests the ReadResponseAsJson function.
func TestReadResponseAsJson(t *testing.T) {

	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Success case: 200 OK with valid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"test", "age":30}`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get response from test server: %v", err)
	}

	var ts testStruct
	if err := ReadResponseAsJson(resp, &ts); err != nil {
		t.Errorf("ReadResponseAsJson failed for success case: %v", err)
	}
	if ts.Name != "test" || ts.Age != 30 {
		t.Errorf("Expected {Name:test, Age:30}, got {Name:%s, Age:%d}", ts.Name, ts.Age)
	}

	// Failure case: 400 Bad Request with error JSON
	serverErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid input"}`))
	}))
	defer serverErr.Close()

	respErr, err := http.Get(serverErr.URL)
	if err != nil {
		t.Fatalf("Failed to get error response from test server: %v", err)
	}

	var tsErr testStruct
	if err := ReadResponseAsJson(respErr, &tsErr); err == nil {
		t.Error("Expected an error for 400 status code, got nil")
	}
	expectedErr := "unexpected status code 400: {\"error\":\"invalid input\"}"
	if err != nil && !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error containing %q, got %q", expectedErr, err.Error())
	}

	// Invalid JSON case: 200 OK with invalid JSON
	serverInvalidJson := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer serverInvalidJson.Close()

	respInvalidJson, err := http.Get(serverInvalidJson.URL)
	if err != nil {
		t.Fatalf("Failed to get invalid JSON response from test server: %v", err)
	}

	var tsInvalid testStruct
	if err := ReadResponseAsJson(respInvalidJson, &tsInvalid); err == nil {
		t.Error("Expected an error for invalid JSON, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "json unmarshal failed") {
		t.Errorf("Expected JSON unmarshal error, got %q", err.Error())
	}
}

// TestCloseResponse tests the CloseResponse function.
func TestCloseResponse(t *testing.T) {
	// Case 1: Close without error, no callback
	resp1 := &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("test"))), Request: &http.Request{URL: &url.URL{Host: "example.com"}}}
	CloseResponse(resp1)

	// Case 2: Close with error, no callback (logs error, but hard to test slog directly)
	// We'll just ensure it doesn't panic and the error is handled internally.
	resp2 := &http.Response{Body: &errorReader{err: errors.New("close error")}, Request: &http.Request{URL: &url.URL{Host: "example.com"}}}
	CloseResponse(resp2)

	// Case 3: Close with error, with callback
	var capturedErr error
	var wg sync.WaitGroup
	wg.Add(1)
	resp3 := &http.Response{Body: &errorReader{err: errors.New("callback error")}, Request: &http.Request{URL: &url.URL{Host: "example.com"}}}
	CloseResponse(resp3, func(err error) {
		capturedErr = err
		wg.Done()
	})
	wg.Wait()
	if capturedErr == nil || capturedErr.Error() != "callback error" {
		t.Errorf("Expected callback to capture 'callback error', got %v", capturedErr)
	}

	// Case 4: Close without error, with callback (callback should not be called)
	callbackCalled := false
	resp4 := &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("test"))), Request: &http.Request{URL: &url.URL{Host: "example.com"}}}
	CloseResponse(resp4, func(err error) {
		callbackCalled = true
	})
	if callbackCalled {
		t.Error("Callback should not be called if Close() returns no error")
	}
}

// errorReader is a mock io.ReadCloser that returns an error on Close().
type errorReader struct {
	io.Reader
	err error
}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF // Simulate end of file
}

func (er *errorReader) Close() error {
	return er.err
}

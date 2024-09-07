package proxy_test

import (
	"github.com/editorpost/spider/collect/proxy"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func validProxyServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("expected content"))
	}))
}

func invalidProxyServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

func TestCheckValidProxy(t *testing.T) {
	proxyServer := validProxyServer()
	defer proxyServer.Close()

	err := proxy.Check(proxyServer.URL, "http://example.com", "expected content", 5*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckInvalidProxy(t *testing.T) {
	proxyServer := invalidProxyServer()
	defer proxyServer.Close()

	err := proxy.Check(proxyServer.URL, "http://example.com", "expected content", 5*time.Second)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCheckTimeout(t *testing.T) {
	proxyServer := validProxyServer()
	defer proxyServer.Close()

	err := proxy.Check(proxyServer.URL, "http://example.com", "expected content", 1*time.Nanosecond)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}

func TestCheckResponseDoesNotContainExpectedString(t *testing.T) {
	proxyServer := validProxyServer()
	defer proxyServer.Close()

	err := proxy.Check(proxyServer.URL, "http://example.com", "unexpected content", 5*time.Second)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCheckInvalidProxyURL(t *testing.T) {
	err := proxy.Check("invalid-url", "http://example.com", "expected content", 5*time.Second)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

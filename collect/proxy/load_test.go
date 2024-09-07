package proxy_test

import (
	"encoding/json"
	"github.com/editorpost/spider/collect/proxy"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadJSONListValidURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxies := []*proxy.Proxy{
			NewProxy("http://proxy1.com"),
			NewProxy("http://proxy2.com"),
		}
		json.NewEncoder(w).Encode(proxies)
	}))
	defer server.Close()

	result := proxy.LoadJSONList(server.URL)
	if len(result) != 2 {
		t.Fatalf("expected 2 proxies, got %d", len(result))
	}
}

func TestLoadJSONListEmptyURL(t *testing.T) {
	result := proxy.LoadJSONList("")
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestLoadJSONListInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	proxy.LoadJSONList(server.URL)
}

func TestLoadJSONListNoProxiesFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]*proxy.Proxy{})
	}))
	defer server.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	proxy.LoadJSONList(server.URL)
}

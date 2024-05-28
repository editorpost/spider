//go:build e2e
// +build e2e

package proxy_test

import (
	"context"
	"fmt"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestNewPool(t *testing.T) {

	p := proxy.NewPool("https://octopart.com/irfb3077pbf-infineon-65873800")
	p.SetCheckContent("Price and Stock")
	require.NoError(t, p.Start())

	u, err := p.GetProxyURL(nil)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	fmt.Println(u.String())
}

type ProxyInfoContextKey struct{}

// ProxySelector returns a dynamically chosen proxy URL based on the request
func ProxySelector(req *http.Request) (*url.URL, error) {
	return url.Parse("http://45.95.203.159:4444") // set actual proxy URL here
}

func TestContext(t *testing.T) {

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			proxyURL, err := ProxySelector(req)
			if err != nil {
				return nil, err
			}
			if proxyURL != nil {
				req = req.WithContext(context.WithValue(req.Context(), ProxyInfoContextKey{}, proxyURL.String()))
			}
			return proxyURL, nil
		},
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Retrieving proxy information from the context, if available
	if proxy, ok := req.Context().Value(ProxyInfoContextKey{}).(string); ok {
		log.Println("Proxy used:", proxy)
	}
}

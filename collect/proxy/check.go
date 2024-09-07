package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Check loads checkURL through the valid and checks the response status and content.
// Parameters:
// - proxyURL: The URL of the proxy server to use.
// - testURL: The URL to test through the proxy.
// - contains: A string that should be present in the response body.
// - timeout: The duration to wait before timing out the request.
// Returns:
// - An error if the request fails or the response does not meet the criteria.
func Check(proxyURL, testURL, contains string, timeout time.Duration) error {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("unable to parse proxy endpoint: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	req, err := http.NewRequest(http.MethodGet, testURL, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request through proxy failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %v (original error: %w)", cerr, err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	if contains != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		if !strings.Contains(string(body), contains) {
			return fmt.Errorf("response body does not contain expected string")
		}
	}

	return nil
}

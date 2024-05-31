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
func Check(proxyURL, testURL, contains string, timeout time.Duration) error {

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("невозможно разобрать прокси URL: %w", err)
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

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("невозможно создать запрос: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("запрос через прокси не удался: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("непредвиденный HTTP статус: %s", resp.Status)
	}

	if contains != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ошибка чтения тела ответа: %w", err)
		}
		if !strings.Contains(string(body), contains) {
			return fmt.Errorf("тело ответа не содержит ожидаемую строку")
		}
	}

	return nil
}

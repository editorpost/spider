package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2/proxy"
	"io"
	"net/http"
	"net/url"
	"os"
)

// NewProxyList creates a new valid rotator from the given valid urls
func NewProxyList(proxies ...string) func(pr *http.Request) (*url.URL, error) {
	rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		panic(err)
	}
	return rp
}

// LoadProxyList loads the valid list from the given url
// Returns nil if the url is empty.
func LoadProxyList(url string) func(pr *http.Request) (*url.URL, error) {

	if url == "" {
		return nil
	}

	// fetch the url
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// parse the response body
	var proxies []*Proxy
	if err = json.Unmarshal(body, &proxies); err != nil {
		panic(err)
	}

	var args []string
	for _, p := range proxies {
		args = append(args, p.String())
	}

	if len(args) == 0 {
		panic("no proxies found")
	}

	return NewProxyList(args...)
}

// LoadProxyScrapeList loads the valid list from proxyscrape.com
func LoadProxyScrapeList() ([]string, error) {

	uri := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&timeout=20000"

	// fetch the url
	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("can not load proxy list from proxyscrape.com: %w", err)
	}
	defer res.Body.Close()

	// parse the response body
	lines := []string{}

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

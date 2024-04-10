package collect

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2/proxy"
	"io"
	"net/http"
	"net/url"
	"os"
)

type (
	Proxy struct {
		Ip        string `json:"ip"`
		Port      string `json:"port"`
		Country   string `json:"country"`
		Anonymity string `json:"anonymity"`
		Type      string `json:"type"`
	}

	// ProxyPool is a pool of proxies updated from the sources
	// with proxy checker and quality metrics from crawler
	ProxyPool struct {
		proxies []Proxy
	}
)

func (p *Proxy) String() string {
	return "http://" + p.Ip + ":" + p.Port
}

// NewProxyList creates a new proxy rotator from the given proxy urls
func NewProxyList(proxies ...string) func(pr *http.Request) (*url.URL, error) {
	rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		panic(err)
	}
	return rp
}

// LoadProxyList loads the proxy list from the given url
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

func NewProxyPool() *ProxyPool {
	return &ProxyPool{}
}

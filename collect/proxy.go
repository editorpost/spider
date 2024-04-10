package collect

import (
	"github.com/gocolly/colly/v2/proxy"
	"net/http"
	"net/url"
)

// ProxyList is a list of proxy urls
type ProxyList []string

// NewProxyList creates a new proxy rotator from the given proxy urls
func NewProxyList(proxies ...string) func(pr *http.Request) (*url.URL, error) {
	rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		panic(err)
	}
	return rp
}

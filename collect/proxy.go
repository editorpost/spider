package collect

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/gocolly/colly/v2"
	"time"
)

// WithProxyPool sets up the proxy for the crawler.
//
// Summary:
// - use proxy.Pool as RoundTripper
// - skip proxyFn if proxy.Pool is used
// - default request timeout is 15 seconds
// - enable cookies (on error - log and skip setting cookies)
// - finally, it sets up retries for response and proxy errors.
func WithProxyPool(args *config.Args) (colly.CollectorOption, error) {

	var (
		err       error
		poolReady bool
		proxies   *proxy.Pool
	)

	if args.ProxyEnabled {
		proxies, err = proxy.StartPool(args.StartURL, args.ProxySources...)
		if err != nil {
			return nil, err
		}
		poolReady = true
	}

	return func(c *colly.Collector) {

		if poolReady {
			// the transport is used to make HTTP requests.
			// with proxy.Pool it is used to rotate the proxy
			// and collect metrics from good/bad proxy responses
			c.WithTransport(proxies.Transport())
		}

		// increase timeouts due proxy rotation
		c.SetRequestTimeout(15 * time.Second)

	}, nil

}

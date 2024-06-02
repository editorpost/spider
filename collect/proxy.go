package collect

import (
	"net/http/cookiejar"
	"time"
)

// withProxy sets up the proxy for the crawler.
//
// Summary:
// - use proxy.Pool as RoundTripper
// - skip proxyFn if proxy.Pool is used
// - default request timeout is 15 seconds
// - enable cookies (on error - log and skip setting cookies)
// - finally, it sets up retries for response and proxy errors.
func (crawler *Crawler) withProxy() {

	// round tripper
	//
	// the transport is used to make HTTP requests.
	// with proxy.Pool it is used to rotate the proxy
	// and collect metrics from good/bad proxy responses
	if crawler.RoundTripper != nil {
		crawler.collect.WithTransport(crawler.RoundTripper)
	}

	// inject proxy func to transport
	//
	// arg: in case of proxy pool, skip this step
	// 	    as the transport is already set and have the proxy func.
	// dev: order matters, call after transport already set
	//      or default transport will be used
	if crawler.ProxyFn != nil {
		crawler.collect.SetProxyFunc(crawler.ProxyFn)
	}

	// timeouts
	crawler.collect.SetRequestTimeout(15 * time.Second)

	// cookie handling
	// for turning off - crawler.collect.DisableCookies()
	j, err := cookiejar.New(&cookiejar.Options{})
	if err == nil {
		crawler.collect.SetCookieJar(j)
	}

	// retries
	crawler.errRetry = NewRetry(ResponseRetries)
	crawler.proxyRetry = NewRetry(BadProxyRetries)
}

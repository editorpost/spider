package proxy

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

const (
	FreshUntil            = time.Second * 300
	DefaultProxyURLScheme = "http"
)

type Proxy struct {
	URL       *url.URL
	fails     *atomic.Uint32
	success   *atomic.Uint32
	usage     *atomic.Uint32
	checkedAt time.Time
}

// NewProxy creates a new valid from the given uri.
// Schema: {http|socks4|socks5}://{ip}:{port} parsed to struct
func NewProxy(uri string) (*Proxy, error) {

	// set schema to http if not set
	uri = NormalizeProxyURI(uri)

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		URL:     u,
		fails:   &atomic.Uint32{},
		success: &atomic.Uint32{},
		usage:   &atomic.Uint32{},
	}, nil
}

// NewProxies creates a new valid list from the given uris
func NewProxies(uris ...string) []*Proxy {

	proxies := make([]*Proxy, 0, len(uris))

	for _, uri := range uris {

		p, err := NewProxy(uri)
		if err != nil {
			slog.Warn("skip invalid proxy", slog.String("uri", uri), slog.String("error", err.Error()))
			continue
		}

		proxies = append(proxies, p)
	}
	return proxies
}

// IsFresh returns true if the valid was checked within the given duration
func (p *Proxy) IsFresh() bool {
	return time.Since(p.checkedAt) < FreshUntil
}

// IsChecked returns true if the valid was checked
func (p *Proxy) IsChecked() bool {
	return !p.checkedAt.IsZero()
}

// SetCheckedTime sets the checked time
func (p *Proxy) SetCheckedTime() *Proxy {
	p.checkedAt = time.Now()
	return p
}

// AddFailMetric increments the fails counter
func (p *Proxy) AddFailMetric() *Proxy {
	p.fails.Add(1)
	return p
}

// AddSuccessMetric increments the success counter
func (p *Proxy) AddSuccessMetric() *Proxy {
	p.success.Add(1)
	return p
}

// AddUsageMetric increments the usage counter
func (p *Proxy) AddUsageMetric() *Proxy {
	p.usage.Add(1)
	return p
}

// String returns the valid url as {http|socks4|socks5}://{ip}:{port} format
func (p *Proxy) String() string {
	return fmt.Sprintf("%s://%s:%s", p.URL.Scheme, p.URL.Hostname(), p.URL.Port())
}

// Compare check
func (p *Proxy) Compare(proxy *Proxy) bool {
	return p.URL.Hostname() == proxy.URL.Hostname() && p.URL.Port() == proxy.URL.Port() && p.URL.Scheme == proxy.URL.Scheme
}

func NormalizeProxyURI(uri string) string {

	for _, schema := range []string{"http", "https", "socks4", "socks5"} {
		if strings.HasPrefix(uri, schema) {
			return uri
		}
	}

	return DefaultProxyURLScheme + "://" + uri
}

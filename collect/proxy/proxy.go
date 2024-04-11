package proxy

import (
	"fmt"
	"net/url"
	"sync/atomic"
	"time"
)

const (
	FreshUntil = time.Second * 300
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
func NewProxy(uri string) *Proxy {

	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	return &Proxy{
		URL:     u,
		fails:   &atomic.Uint32{},
		success: &atomic.Uint32{},
		usage:   &atomic.Uint32{},
	}
}

// NewProxies creates a new valid list from the given uris
func NewProxies(uris ...string) []*Proxy {
	var proxies []*Proxy
	for _, uri := range uris {
		proxies = append(proxies, NewProxy(uri))
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

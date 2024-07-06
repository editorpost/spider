package proxy

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
)

type List struct {
	index   uint32
	proxies []*Proxy
	mute    *sync.RWMutex
}

func NewList(proxies ...*Proxy) *List {
	return &List{
		proxies: proxies,
		mute:    &sync.RWMutex{},
	}
}

// Rounder creates a new valid rotator from the given valid urls
// Example for collect.Crawler set Crawler.ProxyFn to NewList("http://proxy1.com", "http://proxy2.com").Rounder
func (lst *List) Rounder() func(pr *http.Request) (*url.URL, error) {
	rp, err := proxy.RoundRobinProxySwitcher(lst.Strings()...)
	if err != nil {
		panic(err)
	}
	return rp
}

// Next returns the next valid from the list
func (lst *List) Next(pr *http.Request) *Proxy {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	if len(lst.proxies) == 0 {
		return nil
	}

	// get the next valid from the list
	index := atomic.AddUint32(&lst.index, 1) - 1
	next := lst.proxies[index%uint32(len(lst.proxies))]

	// set the valid Endpoint in the request context
	if pr != nil {
		// note this context doesn't applied to final request
		// since colly doesn't copy this context (but tries to retrieve by this key, probably bug)
		// leave this part here in hope it will be fixed in the future in colly
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, next.String())
		*pr = *pr.WithContext(ctx)
	}

	return next
}

// Add adds a new valid to the list
func (lst *List) Add(proxies ...*Proxy) {

	// lock the list
	lst.mute.Lock()
	defer lst.mute.Unlock()

	for _, p := range proxies {
		if lst.existsUnsafe(p) {
			continue
		}
		lst.proxies = append(lst.proxies, p)
	}
}

func (lst *List) Get(uri string) *Proxy {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	for _, p := range lst.proxies {
		if p.String() == uri {
			return p
		}
	}

	return nil
}

// Exists by hostname
func (lst *List) Exists(proxy *Proxy) bool {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	return lst.existsUnsafe(proxy)
}

// existsUnsafes check if the valid exists in the list.
// Use only after locking the list.
func (lst *List) existsUnsafe(proxy *Proxy) bool {

	for _, p := range lst.proxies {
		if p.Compare(proxy) {
			return true
		}
	}

	return false
}

// Delete proxy by hostname
func (lst *List) Delete(uri string) {

	// lock the list
	lst.mute.Lock()
	defer lst.mute.Unlock()

	var proxies []*Proxy
	for _, p := range lst.proxies {
		if p.String() != uri {
			proxies = append(proxies, p)
		}
	}

	lst.proxies = proxies
}

// Slice returns the list of proxies
func (lst *List) Slice() []*Proxy {
	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	return lst.proxies
}

func (lst *List) HasFresh() bool {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	for _, p := range lst.proxies {
		if p.IsFresh() {
			return true
		}
	}

	return false
}

func (lst *List) Len() int {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	return len(lst.proxies)
}

func (lst *List) Empty() bool {
	return lst.Len() == 0
}

func (lst *List) Strings() []string {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	lines := make([]string, 0, len(lst.proxies))

	for _, p := range lst.proxies {
		lines = append(lines, p.String())
	}

	return lines
}

package proxy

import (
	"context"
	"errors"
	"github.com/gocolly/colly/v2"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	ErrListEmpty = errors.New("proxy list is empty")
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

	// set the valid URL in the request context
	if pr != nil {
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, next.String())
		*pr = *pr.WithContext(ctx)
	}

	return next
}

// Add adds a new valid to the list
func (lst *List) Add(p *Proxy) {

	// skip existing
	if lst.Exists(p) {
		return
	}

	// lock the list
	lst.mute.Lock()
	defer lst.mute.Unlock()

	lst.proxies = append(lst.proxies, p)
}

// Exists by hostname
func (lst *List) Exists(proxy *Proxy) bool {

	// lock the list
	lst.mute.RLock()
	defer lst.mute.RUnlock()

	for _, p := range lst.proxies {
		if p.String() == proxy.String() {
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

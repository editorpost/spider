package proxy_test

import (
	"github.com/editorpost/spider/collect/proxy"
	"net/url"
	"strconv"
	"sync"
	"testing"
)

func TestListConcurrency(_ *testing.T) {

	proxies := []*proxy.Proxy{
		NewProxy("http://proxy1.com"),
		NewProxy("http://proxy2.com"),
		NewProxy("http://proxy3.com"),
	}

	lst := proxy.NewList(proxies...)

	var wg sync.WaitGroup

	// Concurrently add proxies
	for i := 4; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lst.Add(NewProxy("http://proxy" + strconv.Itoa(i) + ".com"))
		}(i)
	}

	// Concurrently get proxies
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lst.Next(nil)
		}()
	}

	// Concurrently check existence
	for i := 4; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lst.Exists(NewProxy("http://proxy" + strconv.Itoa(i) + ".com"))
		}(i)
	}

	wg.Wait()
}

func TestAddProxy(t *testing.T) {
	lst := proxy.NewList()
	proxy1 := NewProxy("http://proxy1.com")
	lst.Add(proxy1)

	if !lst.Exists(proxy1) {
		t.Fatalf("expected proxy to exist in the list")
	}
}

func TestAddDuplicateProxy(t *testing.T) {
	lst := proxy.NewList()
	proxy1 := NewProxy("http://proxy1.com")
	lst.Add(proxy1, proxy1)

	if len(lst.Slice()) != 1 {
		t.Fatalf("expected only one instance of the proxy in the list")
	}
}

func TestDeleteProxy(t *testing.T) {
	lst := proxy.NewList()
	proxy1 := NewProxy("http://proxy1.com")
	lst.Add(proxy1)
	lst.Delete(proxy1.String())

	if lst.Exists(proxy1) {
		t.Fatalf("expected proxy to be deleted from the list")
	}
}

func TestNextProxy(t *testing.T) {
	lst := proxy.NewList(NewProxy("http://proxy1.com"), NewProxy("http://proxy2.com"))
	proxy1 := lst.Next(nil)
	proxy2 := lst.Next(nil)

	if proxy1 == nil || proxy2 == nil {
		t.Fatalf("expected to get proxies from the list")
	}

	if proxy1 == proxy2 {
		t.Fatalf("expected to get different proxies in round-robin fashion")
	}
}

func TestEmptyList(t *testing.T) {
	lst := proxy.NewList()

	if !lst.Empty() {
		t.Fatalf("expected list to be empty")
	}
}

func TestHasFreshProxy(t *testing.T) {
	lst := proxy.NewList(NewProxy("http://proxy1.com"))
	if lst.HasFresh() {
		t.Fatalf("expected no fresh proxies in the list")
	}
}

func NewProxy(addr string) *proxy.Proxy {
	return &proxy.Proxy{URL: &url.URL{Host: addr}}
}

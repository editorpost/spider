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

func NewProxy(addr string) *proxy.Proxy {
	return &proxy.Proxy{URL: &url.URL{Host: addr}}
}

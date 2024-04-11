package proxy

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Pool is a pool of proxies checked against the test URL.
// Periodically checks the proxies and updates the valid list.
type Pool struct {
	valid        *List
	check        *List
	checkURL     string
	checkContent string
	checkTimeout time.Duration
	// Loader is a function to load the proxy list
	Loader func() ([]string, error)
	// Checker is a function to check the proxy by URI string
	Checker func(string) error
	mute    sync.RWMutex
}

func NewPool(testURL string) *Pool {
	return &Pool{
		valid:        NewList(),
		check:        NewList(),
		checkURL:     testURL,
		checkTimeout: time.Second * 30,
		Checker:      nil,
	}
}

// Start initializes the pool with the given proxies.
func (pool *Pool) Start() error {

	if err := pool.load(); err != nil {
		return err
	}

	go pool.Check()

	// log metrics every 30 seconds
	go pool.Report()

	return nil
}

// GetProxyURL returns the next valid from the pool or blocks until one is available.
// Every 30 seconds prints report of the valid pool.
func (pool *Pool) GetProxyURL(pr *http.Request) (*url.URL, error) {

	// load next valid proxy
	if proxy := pool.valid.Next(pr); proxy != nil {
		slog.Info("with proxy", slog.String("url", proxy.URL.String()))
		proxy.AddUsageMetric()
		return proxy.URL, nil
	}

	// wait for a valid proxy
	start := time.Now()
	try := 0
	reportEvery := time.Second * 30
	reportAt := time.Time{}

	// run 12 hours since start
	for time.Since(start) < time.Hour*12 {

		try++
		time.Sleep(time.Second)

		// try
		if proxy := pool.valid.Next(pr); proxy != nil {
			proxy.AddUsageMetric()
			return proxy.URL, nil
		}

		// report
		if time.Now().After(reportAt) {
			slog.Info("waiting valid proxies", slog.Duration("elapsed", time.Since(start)), slog.Int("try", try))
			reportAt = time.Now().Add(reportEvery)
		}
	}

	// no valid proxies after 12 hours
	slog.Error("no valid proxies", slog.Duration("elapsed", time.Since(start)), slog.Int("try", try))

	return nil, errors.New("no valid proxies after 12 hours")
}

func (pool *Pool) load() error {

	if pool.Loader == nil {
		pool.Loader = LoadPublicLists
	}

	proxies, err := pool.Loader()
	if err != nil {
		return err
	}

	for _, proxy := range proxies {
		pool.check.Add(NewProxy(proxy))
	}

	return nil
}

func (pool *Pool) Check() {
	var wg sync.WaitGroup
	for _, p := range pool.check.Slice() {
		wg.Add(1)
		go pool.CheckProxy(p, &wg)
	}
	wg.Wait()
}

func (pool *Pool) CheckProxy(p *Proxy, wg *sync.WaitGroup) {

	defer wg.Done()
	p.SetCheckedTime()

	// checker with fallback
	checker := pool.Checker
	if checker == nil {
		checker = func(proxyURL string) error {
			return Check(proxyURL, pool.checkURL, pool.checkContent, pool.checkTimeout)
		}
	}

	// check proxy
	if err := checker(p.String()); err != nil {
		// failed proxy
		p.AddFailMetric()
		pool.check.Delete(p.String())
		return
	}

	// success proxy
	p.AddSuccessMetric()
	pool.valid.Add(p)
}

func (pool *Pool) SetCheckContent(contains string) {
	pool.checkContent = contains
}

func (pool *Pool) SetCheckTimeout(timeout time.Duration) {
	pool.checkTimeout = timeout
}

func (pool *Pool) Report() {
	for {
		time.Sleep(time.Second * 30)
		slog.Info("valid proxies", slog.Int("count", len(pool.valid.Slice())))
	}
}

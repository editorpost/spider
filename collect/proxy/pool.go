package proxy

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

var (
	ErrBadProxy = errors.New("bad proxy")
)

// Pool is a pool of proxies checked against the test URL.
// Periodically checks the proxies and updates the valid list.
type Pool struct {
	valid        *List
	check        *List
	suspect      *List
	checkURL     string
	checkContent string
	checkTimeout time.Duration
	// Loader is a function to load the proxy list
	Loader func() ([]string, error)
	// Checker is a function to check the proxy by URI string
	Checker func(string) error
	mute    sync.RWMutex
	rtp     *http.Transport
}

func NewPool(testURL string) *Pool {

	pool := &Pool{
		valid:        NewList(),
		check:        NewList(),
		suspect:      NewList(),
		checkURL:     testURL,
		checkTimeout: time.Second * 30,
		Checker:      nil,
	}

	pool.rtp = &http.Transport{
		Proxy:             pool.GetProxyURL,
		DisableKeepAlives: true,
		OnProxyConnectResponse: func(ctx context.Context, proxyURL *url.URL, req *http.Request, resp *http.Response) error {

			p := pool.valid.Get(proxyURL.String())
			p.AddUsageMetric()

			if resp.StatusCode == http.StatusOK {
				p.AddSuccessMetric()
				return nil
			}

			if p.fails.Load() > 10 {
				pool.valid.Delete(proxyURL.String())
				pool.suspect.Add(p)
				return nil
			}

			p.AddFailMetric()

			return ErrBadProxy
		},
	}

	return pool
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

func (pool *Pool) Transport() *http.Transport {
	return pool.rtp
}

// GetProxyURL returns the next valid from the pool or blocks until one is available.
// Every 30 seconds prints report of the valid pool.
func (pool *Pool) GetProxyURL(pr *http.Request) (*url.URL, error) {

	// load next valid proxy
	if proxy := pool.valid.Next(pr); proxy != nil {
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

	proxies, loadErr := pool.Loader()
	if loadErr != nil {
		return loadErr
	}

	for _, proxy := range proxies {
		p, err := NewProxy(proxy)
		if err != nil {
			slog.Warn("skip invalid proxy", slog.String("uri", proxy), slog.String("error", err.Error()))
			continue
		}
		pool.check.Add(p)
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
		pool.check.Delete(p.String())
		return
	}

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
		time.Sleep(time.Second * 60)
		slog.Info("valid proxies", slog.Int("count", len(pool.valid.Slice())))

		pool.ReportFiveMostFailed()
		pool.ReportFiveMostUsed()
	}
}

func (pool *Pool) ReportFiveMostFailed() {

	proxies := pool.valid.Slice()
	if len(proxies) == 0 {
		return
	}

	// sort by fails
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].fails.Load() > proxies[j].fails.Load()
	})

	slog.Info("MOST FAILED PROXIES")

	for _, p := range proxies[:min(5, len(proxies))] {
		slog.Info("proxy",
			slog.Int("success", int(p.success.Load())),
			slog.Int("fails", int(p.fails.Load())),
			slog.Int("usage", int(p.usage.Load())),
			slog.String("url", p.URL.String()),
		)
	}
}

func (pool *Pool) ReportFiveMostUsed() {

	proxies := pool.valid.Slice()
	if len(proxies) == 0 {
		return
	}

	// sort by usage
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].usage.Load() > proxies[j].usage.Load()
	})
	slog.Info("TOP USED PROXIES")

	for _, p := range proxies[:min(5, len(proxies))] {
		slog.Info("proxy",
			slog.Int("success", int(p.success.Load())),
			slog.Int("fails", int(p.fails.Load())),
			slog.Int("usage", int(p.usage.Load())),
			slog.String("url", p.URL.String()),
		)
	}
}

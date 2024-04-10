package collect

import (
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/http"
	"sync"
)

const (
	// RequestRetries is the number of retries for request errors from colly.onError handler.
	// it handler network errors, timeouts, etc.
	RequestRetries = 2
	// ResponseRetries is the number of retries for response errors from crawler.visit handler.
	// it handler http status codes, anti-scraping, captcha, etc.
	ResponseRetries = 15
)

type Retry struct {
	_count map[string]uint16
	mute   *sync.Mutex
}

func NewRetry() *Retry {
	return &Retry{
		_count: make(map[string]uint16),
		mute:   &sync.Mutex{},
	}
}

func (r *Retry) Request(resp *colly.Response) bool {

	if r.Limited(resp) {
		return false
	}

	r.inc(resp.Request.URL.String())
	err := resp.Request.Retry()

	// not actual response error, since request might be executed in async mode
	if err != nil {
		slog.Error("Request failed",
			slog.String("url", resp.Request.URL.String()),
			slog.String("proxy", resp.Request.ProxyURL),
			slog.String("err", err.Error()),
		)
	}

	return true
}

func (r *Retry) Limited(resp *colly.Response) bool {

	// skip some errors
	count := r.Count(resp.Request.URL.String())

	if resp.StatusCode == http.StatusForbidden {
		if count > ResponseRetries {
			return false
		}
	}

	return count > RequestRetries
}

func (r *Retry) Count(url string) uint16 {

	r.mute.Lock()
	defer r.mute.Unlock()

	if _, ok := r._count[url]; !ok {
		r._count[url] = 0
	}

	return r._count[url]
}

func (r *Retry) inc(url string) {
	r.mute.Lock()
	defer r.mute.Unlock()
	r._count[url]++
}
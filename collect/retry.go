package collect

import (
	"github.com/gocolly/colly/v2"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	// BadProxyRetries is the number of retries for request errors from colly.onError handler.
	// it handler network errors, timeouts, etc.
	BadProxyRetries = 15
	// ResponseRetries is the number of retries for response errors from crawler.visit handler.
	// it handler http status codes, anti-scraping, captcha, etc.
	ResponseRetries = 3
)

type Retry struct {
	limit  uint16
	mute   *sync.Mutex
	_count map[string]uint16
}

func NewRetry(limit uint16) *Retry {
	return &Retry{
		limit:  limit,
		mute:   &sync.Mutex{},
		_count: make(map[string]uint16),
	}
}

func (r *Retry) Request(resp *colly.Response) bool {

	if r.Limited(resp) {
		return false
	}

	r.inc(resp.Request.URL.String())

	tries := &atomic.Int32{}

	for {

		tries.Add(1)
		err := resp.Request.Retry()

		// not actual response error, since request might be executed in async mode
		if err == nil {
			return true
		}

		if tries.Load() >= int32(r.limit) {
			slog.Debug("retry failed",
				slog.String("url", resp.Request.URL.String()),
				slog.String("err", err.Error()),
			)
			// limit reached
			return false
		}
	}

}

func (r *Retry) Limited(resp *colly.Response) bool {

	url := resp.Request.URL.String()
	count := r.Count(url)

	return count > r.limit
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

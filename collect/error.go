package collect

import (
	"errors"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
)

// error logging
func (crawler *Crawler) error(resp *colly.Response, err error) {

	if errors.Is(err, proxy.ErrBadProxy) {

		// retry on error with new proxy candidate
		if crawler.proxyRetry.Request(resp) {
			return
		}

		slog.Error("cannot find proper proxy candidate for the url",
			slog.String("url", resp.Request.URL.String()),
			slog.String("proxy", resp.Request.ProxyURL),
			slog.Int("status", resp.StatusCode),
		)

		return
	}

	// catch *url.Error
	var urlErr *url.Error
	if errors.As(err, &urlErr) {

		// retry on error with new working proxy
		if crawler.errRetry.Request(resp) {
			return
		}

		slog.Debug("url error",
			slog.String("err", err.Error()),
			slog.String("url", resp.Request.URL.String()),
			slog.String("proxy", resp.Request.ProxyURL),
			slog.Int("status", resp.StatusCode),
		)
		return
	}

	slog.Debug("response failed",
		slog.String("err", err.Error()),
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.Int("status", resp.StatusCode),
	)
}

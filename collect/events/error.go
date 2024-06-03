package events

import (
	"errors"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
)

// error logging
func (crawler *Dispatch) error(resp *colly.Response, err error) {

	crawler.deps.Monitor.OnError(resp, err)

	if errors.Is(err, proxy.ErrBadProxy) {

		// retry on error with new proxy candidate
		if crawler.proxyRetry.Request(resp) {
			return
		}

		slog.Debug("bad proxy",
			slog.String("url", resp.Request.URL.String()),
			slog.String("proxy", resp.Request.ProxyURL),
			slog.Int("status", resp.StatusCode),
		)

		return
	}

	// catch *url.OnError
	var urlErr *url.Error
	if errors.As(err, &urlErr) {

		// retry on error with new working proxy
		if crawler.errorRetry.Request(resp) {
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

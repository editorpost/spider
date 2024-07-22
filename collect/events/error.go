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

		LogRespError("bad proxy", resp, err)
		return
	}

	// catch *url.OnError
	var urlErr *url.Error
	if errors.As(err, &urlErr) {

		// retry on error with new working proxy
		if crawler.errorRetry.Request(resp) {
			return
		}

		LogRespError("url error", resp, err)
		return
	}

	LogRespError("response error", resp, err)
}

func LogRespError(msg string, resp *colly.Response, err error) {
	slog.Debug(msg,
		slog.String("err", err.Error()),
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.Int("status", resp.StatusCode),
	)
}

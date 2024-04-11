package collect

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
)

// error logging
func (crawler *Crawler) error(resp *colly.Response, err error) {

	// slog.Error(err.Error())

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = urlErr.Err
		if crawler.retry.Request(resp) {
			return
		}
	}

	// error if retry limit is reached
	slog.Error(err.Error(),
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.Int("status", resp.StatusCode),
	)
}

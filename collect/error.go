package collect

import (
	"github.com/gocolly/colly/v2"
	"log/slog"
)

// error logging
func (crawler *Crawler) error(resp *colly.Response, err error) {

	if crawler.retry.Request(resp) {
		return
	}

	// error if retry limit is reached
	slog.Error(err.Error(),
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.Int("status", resp.StatusCode),
	)
}

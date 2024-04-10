package collect

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
)

// error logging
func (crawler *Crawler) error(resp *colly.Response, err error) {

	// skip some errors
	if crawler.ignoreError(err) {
		return
	}

	// write error log
	slog.Warn("visit failed",
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.String("err", err.Error()),
	)

	slog.Error("crawler failed",
		slog.String("error", err.Error()),
		slog.String("url", resp.Request.URL.String()),
		slog.String("proxy", resp.Request.ProxyURL),
		slog.String("query", crawler.EntitySelector),
	)
}

func (crawler *Crawler) ignoreError(err error) bool {
	skipErrors := []error{
		colly.ErrAlreadyVisited,
		colly.ErrForbiddenDomain,
		colly.ErrForbiddenURL,
		colly.ErrNoURLFiltersMatch,
	}

	for _, skip := range skipErrors {
		if errors.Is(err, skip) {
			return true
		}
	}
	return false
}

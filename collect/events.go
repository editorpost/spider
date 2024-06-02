package collect

import (
	"errors"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
	"strings"
)

// withEventHandlers sets up the event handlers for the crawler.
// It sets handlers for HTML elements, errors, requests, and responses.
func (crawler *Crawler) withEventHandlers() {
	// collect links
	crawler.collect.OnHTML(`a[href]`, crawler.visit())
	// extract data
	crawler.collect.OnHTML(`html`, crawler.extract())
	// catch errors, run retry
	crawler.collect.OnError(crawler.error)
	// rest for monitoring
	crawler.collect.OnRequest(crawler.request)
	crawler.collect.OnResponse(crawler.response)
	crawler.collect.OnScraped(crawler.scraped)
}

// visit links found in the DOM
func (crawler *Crawler) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		// absolute url
		link := e.Request.AbsoluteURL(e.Attr("href"))

		// skip empty and anchor links
		if link == "" || strings.HasPrefix(link, "#") {
			return
		}

		// skip images, scripts, etc.
		if !isValidURLExtension(link) {
			return
		}

		// visit the link
		if err := crawler.queue.AddURL(link); err != nil {
			slog.Warn("crawler queue", slog.String("error", err.Error()))
		}
	}
}

// extract entries from html selections
func (crawler *Crawler) extract() func(e *colly.HTMLElement) {
	return func(doc *colly.HTMLElement) {

		// entity url regex
		if crawler._entityURL != nil {
			if !crawler._entityURL.MatchString(doc.Request.URL.String()) {
				return
			}
		}

		// selected html selections matching the query
		// might be empty if the query is not found
		for _, selected := range crawler.selections(doc) {

			if err := crawler.Extractor(doc, selected); err != nil {
				crawler.Monitor.OnError(doc.Response, err)
				continue
			}

			crawler.Monitor.OnExtract(doc.Response)
		}
	}
}

// error logging
func (crawler *Crawler) error(resp *colly.Response, err error) {

	crawler.Monitor.OnError(resp, err)

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

// request dispatcher
func (crawler *Crawler) request(r *colly.Request) {
	crawler.Monitor.OnRequest(r)
}

// response dispatcher
func (crawler *Crawler) response(r *colly.Response) {
	crawler.Monitor.OnResponse(r)
}

// scraped dispatcher
func (crawler *Crawler) scraped(r *colly.Response) {
	crawler.Monitor.OnScraped(r)
}

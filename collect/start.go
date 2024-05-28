package collect

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/gocolly/colly/v2/storage"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"sync/atomic"
)

// Crawler for scraping a website
type Crawler struct {
	// StartURL is the url to start the scraping
	StartURL string
	// AllowedURL is comma separated regex to match the urls
	// use it to reduce the number of urls to visit
	AllowedURL string
	// EntityURL is the regex to match the entity urls
	// use it to extract the entity urls
	EntityURL string
	// EntitySelector is the css selector to match the elements
	// use selector for extracting entities and filtering pages
	EntitySelector string
	// UseBrowser is a flag to use browser for rendering the page
	UseBrowser bool
	// Depth if is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	Depth int
	// UserAgent is the user agent string used by the setup
	UserAgent string
	// ProxyFn is the function to return the proxy for the request
	ProxyFn func(*http.Request) (*url.URL, error)
	// RoundTripper is the function to return the next proxy from the list
	RoundTripper http.RoundTripper
	// Extractor is the function to process matched the data, e.g. html tag node
	Extractor func(*colly.HTMLElement, *goquery.Selection) error
	// Storage is the storage backend for the setup
	Storage storage.Storage

	// jsLoadSuccess count the number of successful fallback JS loads
	jsFallbackSuccess *atomic.Int32
	jsSuccessRequired int32
	collect           *colly.Collector
	_entityURL        *regexp.Regexp
	chromeCtx         context.Context
	errRetry          *Retry
	proxyRetry        *Retry
	report            *Report
	queue             *queue.Queue
}

// Start the scraping Crawler.
func (crawler *Crawler) Start() error {

	slog.Info("start",
		slog.String("url", crawler.StartURL),
		slog.String("allowed", crawler.AllowedURL),
		slog.String("entity", crawler.EntityURL),
		slog.String("selector", crawler.EntitySelector),
		slog.Bool("browser", crawler.UseBrowser),
		slog.Int("depth", crawler.Depth),
	)

	crawler.setup()

	if crawler.UseBrowser {
		// create chrome allocator context
		cancel := crawler.setupChrome()
		// disable async in browser mode
		crawler.collect.Async = false
		defer cancel()
	}

	if err := crawler.queue.AddURL(crawler.StartURL); err != nil {
		return err
	}

	if err := crawler.queue.Run(crawler.collect); err != nil {
		return err
	}

	crawler.collect.Wait()

	return nil
}

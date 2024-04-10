package collect

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/storage"
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
	// UserAgent is the user agent string used by the collector
	UserAgent string
	// ProxyFn is the function to return the proxy for the request
	ProxyFn func(*http.Request) (*url.URL, error)
	// Extractor is the function to process matched the data, e.g. html tag node
	Extractor func(*colly.HTMLElement, *goquery.Selection) error
	// Collector is the storage backend for the collector
	Collector storage.Storage

	// jsLoadSuccess count the number of successful fallback JS loads
	jsFallbackSuccess *atomic.Int32
	jsSuccessRequired int32
	collect           *colly.Collector
	_entityURL        *regexp.Regexp
	chromeCtx         context.Context
	retry             *Retry
}

// Start the scraping Crawler.
func (crawler *Crawler) Start() error {

	collector := crawler.collector()

	if crawler.UseBrowser {
		// create chrome allocator context
		cancel := crawler.setupChrome()
		// disable async in browser mode
		collector.Async = false
		defer cancel()
	}

	if err := collector.Visit(crawler.StartURL); err != nil {
		return err
	}

	collector.Wait()

	return nil
}

package collect

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/storage"
	"log/slog"
	"net/url"
)

// Crawler for scraping a website
type Crawler struct {
	// StartURL is the url to start the scraping
	StartURL string
	// MatchURL is comma separated regex to match the urls
	// use it to reduce the number of urls to visit
	MatchURL string
	// Query is the css selector to match the elements
	// use selector for extracting entities and filtering pages
	Query string
	// Depth if is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	Depth int
	// UserAgent is the user agent string used by the collector
	UserAgent string
	// Extractor is the function to process matched the data, e.g. html tag node
	Extractor func(*colly.HTMLElement, *goquery.Selection) error
	// Collector is the storage backend for the collector
	Collector storage.Storage

	collect *colly.Collector
}

// Start the scraping Crawler.
func (crawler *Crawler) Start() error {

	collector := crawler.collector()

	if err := collector.Visit(crawler.StartURL); err != nil {
		return err
	}

	collector.Wait()

	return nil
}

// error logging
func (crawler *Crawler) error(url string, err error) {
	slog.Error("crawler failed",
		slog.String("error", err.Error()),
		slog.String("url", url),
		slog.String("query", crawler.Query),
	)
}

// MustHost from url
func MustHost(fromURL string) string {

	uri, err := url.Parse(fromURL)
	if err != nil {
		panic(err)
	}

	return uri.Host
}

package events

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
)

type (
	Dispatch struct {
		args       *config.Args
		deps       *config.Deps
		queue      Queue
		browser    Browser
		proxyRetry *Retry
		errorRetry *Retry

		// extractor is the function to process matched the data, e.g. html tag node
		extractor func(*colly.HTMLElement, *goquery.Selection) error
		// Metrics is the spider event dispatcher and VictoriaMetrics
		monitor config.Metrics
	}

	Browser interface {
		Browse(uri string) (*goquery.Selection, error)
	}

	Queue interface {
		AddURL(uri string) error
	}
)

func NewDispatcher(args *config.Args, deps *config.Deps, queue Queue, browser Browser) *Dispatch { // long miles away...
	return &Dispatch{
		args:       args,
		deps:       deps,
		queue:      queue,
		browser:    browser,
		proxyRetry: NewRetry(BadProxyRetries),
		errorRetry: NewRetry(ResponseRetries),
	}
}

// WithDispatcher sets up the event handlers for the crawler.
// It sets handlers for HTML elements, errors, requests, and responses.
// noinspection GoUnusedExportedFunction
func WithDispatcher(args *config.Args, deps *config.Deps, queue Queue, browser Browser) func(*colly.Collector) {

	d := NewDispatcher(args, deps, queue, browser)

	return func(c *colly.Collector) {

		// collect links
		c.OnHTML(`a[href]`, d.visit())
		// extract data
		c.OnHTML(`html`, d.extract())
		// catch errors, run retry
		c.OnError(d.error)
		// rest for monitoring
		c.OnRequest(d.request)
		c.OnResponse(d.response)
		c.OnScraped(d.scraped)
	}
}

// request dispatcher
func (crawler *Dispatch) request(r *colly.Request) {
	crawler.deps.Monitor.OnRequest(r)
}

// response dispatcher
func (crawler *Dispatch) response(r *colly.Response) {
	crawler.deps.Monitor.OnResponse(r)
}

// scraped dispatcher
func (crawler *Dispatch) scraped(r *colly.Response) {
	crawler.deps.Monitor.OnScraped(r)
}

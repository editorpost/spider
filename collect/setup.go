package collect

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/events"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/queue"
	"log/slog"
	"net/http/cookiejar"
	"regexp"
)

// collector based on colly
func (crawler *Crawler) collector() (*colly.Collector, error) {

	// return if already initialized
	if crawler.collect != nil {
		return crawler.collect, nil
	}

	slog.Info("setup collector", crawler.args.Log())

	if err := crawler.withQueue(); err != nil {
		return nil, err
	}

	withProxyPool, err := WithProxyPool(crawler.args)
	if err != nil {
		return nil, err
	}

	// Set up a new collector with a maximum depth and maximum body size
	crawler.collect = colly.NewCollector(
		colly.MaxDepth(crawler.args.Depth),
		colly.MaxBodySize(10<<20), // 10MB
		crawler.VisitUrlsFilter(crawler.args),
		events.WithDispatcher(crawler.args, crawler.deps, events.Queue(crawler.queue), events.Browser(crawler)),
		withProxyPool,
	)

	// revisit the same URL
	crawler.collect.AllowURLRevisit = !crawler.args.VisitOnce

	if err = crawler.collect.SetStorage(crawler.deps.Storage); err != nil {
		return nil, err
	}

	// Set a random user agent
	extensions.RandomUserAgent(crawler.collect)

	// cookie handling
	// for turning off - crawler.collect.DisableCookies()
	j, err := cookiejar.New(&cookiejar.Options{})
	if err == nil {
		crawler.collect.SetCookieJar(j)
	}

	return crawler.collect, nil
}

// VisitUrlsFilter sets up the Endpoint filters for the collector.
// It applies a regular expression filter to the URLs visited by the collector.
// Allowed Endpoint pattern is used to extract links in hope to find entity URLs.
// In other hand, ExtractURL is used to run extractors on the page.
func (crawler *Crawler) VisitUrlsFilter(args *config.Config) colly.CollectorOption {
	return func(collector *colly.Collector) {

		// Append the host of the start Endpoint to the allowed domains of the collector
		collector.AllowedDomains = append(collector.AllowedDomains, config.MustHostname(args.StartURL))

		// Append the allowed Endpoint to the Endpoint filters of the collector
		for _, allowedURL := range args.AllowedURLs {
			crawler.VisitUrlFilter(allowedURL, collector)
		}
	}
}

func (crawler *Crawler) VisitUrlFilter(allowedURL string, collector *colly.Collector) {

	// Generate regular expressions from the start, allowed, and entity URLs
	allowed := config.RegexPattern(allowedURL)

	// URLFilters is a list of regular expressions of allowed urls.
	// if ANY expression matches the URL, the URL is allowed to be visited.
	collector.URLFilters = append(collector.URLFilters, regexp.MustCompile(allowed))
}

// withQueue sets up the request queue for the crawler.
// It creates a new request queue with 25 consumer threads and an in-memory queue storage with a maximum size of 50MB.
// If an error occurs during the collector, it panics and stops the execution.
//
// create a request queue with number of consumer threads
// https://go-colly.org/docs/examples/queue/
func (crawler *Crawler) withQueue() (err error) {

	crawler.queue, err = queue.New(
		5, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 5000000}, // 5MB
	)

	return err
}

//goland:noinspection GoLinter
func (crawler *Crawler) withDebug() {

	if crawler.deps.Debugger == nil {
		return
	}

	// colly event dispatcher for logging, monitoring and debugging
	crawler.collect.SetDebugger(crawler.deps.Debugger)
}

package spider

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
	"regexp"
	"time"
)

type Task struct {
	// StartURL is the url to start the scraping
	StartURL string
	// MatchURL is comma separated regex to match the urls
	MatchURL string
	// Query is the css selector to match the elements
	Query string
	// Depth if is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	Depth int
	// UserAgent is the user agent string used by the collector
	UserAgent string
	// Extract is the function to process matched the data
	Extract func(*colly.HTMLElement) error
}

func Run(cfg Task) error {

	c := Collector(cfg)

	if err := c.Visit(cfg.StartURL); err != nil {
		return err
	}

	c.Wait()

	return nil
}

func Collector(task Task) *colly.Collector {

	// default collector
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(DomainName(task.StartURL), "127.0.0.1"),
		colly.UserAgent(task.UserAgent),
		colly.MaxDepth(task.Depth),
		colly.URLFilters(
			regexp.MustCompile(task.MatchURL),
		),
	)

	// limit parallelism per domain
	rule := &colly.LimitRule{DomainGlob: DomainName(task.StartURL), Parallelism: 2, RandomDelay: time.Second}
	if err := c.Limit(rule); err != nil {
		panic(err)
	}

	// follow links
	c.OnHTML(`a[href]`, OnLink(c))

	// extract data
	c.OnHTML(task.Query, func(e *colly.HTMLElement) {

		// do something
		if err := task.Extract(e); err != nil {
			slog.Error("extract failed",
				slog.String("error", err.Error()),
				slog.String("url", e.Request.URL.String()),
				slog.String("query", task.Query),
			)
		}
	})

	return c
}

// OnLink returns a function that visits the link and ask Plugin if it must be visited.
func OnLink(c *colly.Collector) func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))

		err := c.Visit(link)
		if err == nil {
			return
		}

		if errors.Is(err, colly.ErrAlreadyVisited) {
			return
		}

		if errors.Is(err, colly.ErrForbiddenDomain) {
			return
		}

		slog.Warn("visit failed",
			slog.String("url", link),
			slog.String("err", err.Error()),
		)
	}
}

// DomainName returns the domain from a given url or panic
func DomainName(u string) string {

	uri, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	return uri.Host
}

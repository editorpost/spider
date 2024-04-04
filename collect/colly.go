package collect

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"regexp"
	"time"
)

// collector based on colly
func (crawler *Crawler) collector() *colly.Collector {

	if crawler.collect != nil {
		return crawler.collect
	}

	// default collector
	crawler.collect = colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(MustHost(crawler.StartURL), "127.0.0.1"),
		colly.UserAgent(crawler.UserAgent),
		colly.MaxDepth(crawler.Depth),
		colly.URLFilters(
			regexp.MustCompile(crawler.MatchURL),
		),
	)

	// storage backend
	if crawler.Collector != nil {
		err := crawler.collect.SetStorage(crawler.Collector)
		if err != nil {
			panic(err)
		}
	}

	// limit parallelism per domain
	rule := &colly.LimitRule{DomainGlob: MustHost(crawler.StartURL), Parallelism: 2, RandomDelay: time.Second}
	if err := crawler.collect.Limit(rule); err != nil {
		panic(err)
	}

	crawler.collect.OnHTML(`a[href]`, crawler.visit())
	crawler.collect.OnHTML(`html`, crawler.extract())

	return crawler.collect
}

// visit links found in the DOM
func (crawler *Crawler) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))

		err := crawler.collector().Visit(link)
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

// extract entries from html nodes
func (crawler *Crawler) extract() func(e *colly.HTMLElement) {
	return func(doc *colly.HTMLElement) {

		// selected html nodes matching the query
		// might be empty if the query is not found
		for _, selected := range crawler.nodes(doc) {
			err := crawler.Extractor(doc, selected)
			if err != nil {
				crawler.error(doc.Request.URL.String(), err)
				// explicitly
				continue
			}
		}
	}
}

package collect

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"regexp"
	"time"
)

// collector based on colly
func (task Task) collector() *colly.Collector {

	if task.collect != nil {
		return task.collect
	}

	// default collector
	task.collect = colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(MustHost(task.StartURL), "127.0.0.1"),
		colly.UserAgent(task.UserAgent),
		colly.MaxDepth(task.Depth),
		colly.URLFilters(
			regexp.MustCompile(task.MatchURL),
		),
	)

	// storage backend
	if task.Storage != nil {
		err := task.collect.SetStorage(task.Storage)
		if err != nil {
			panic(err)
		}
	}

	// limit parallelism per domain
	rule := &colly.LimitRule{DomainGlob: MustHost(task.StartURL), Parallelism: 2, RandomDelay: time.Second}
	if err := task.collect.Limit(rule); err != nil {
		panic(err)
	}

	task.collect.OnHTML(`a[href]`, task.visit())
	task.collect.OnHTML(`html`, task.extract())

	return task.collect
}

// visit links found in the DOM
func (task Task) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))

		err := task.collector().Visit(link)
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
func (task Task) extract() func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {

		// selected html nodes matching the query
		// might be empty if the query is not found
		for _, selected := range task.nodes(e) {
			err := task.Extract(selected, e.Request.URL)
			if err != nil {
				task.error(e.Request.URL.String(), err)
				// explicitly
				continue
			}
		}
	}
}

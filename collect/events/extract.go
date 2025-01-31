package events

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
	"regexp"
)

// extract entries from html selections
func (crawler *Dispatch) extract() func(e *colly.HTMLElement) {

	var patterns []*regexp.Regexp

	for _, expr := range crawler.args.ExtractURLs {
		patterns = append(patterns, regexp.MustCompile(config.RegexPattern(expr)))
	}

	match := func(u *url.URL) bool {

		if len(patterns) == 0 {
			return true
		}

		for _, pattern := range patterns {
			if pattern.MatchString(u.String()) {
				return true
			}
		}

		return false
	}

	return func(doc *colly.HTMLElement) {

		// the url matches the expression
		if !match(doc.Request.URL) {
			slog.Info("extract: url not matched",
				slog.String("url", doc.Request.URL.String()),
				slog.String("title", doc.DOM.Find("title").Text()),
			)
			return
		}

		// check extraction limit
		if crawler.IsExtractionLimitReached() {
			slog.Info("extract: limit reached",
				slog.String("url", doc.Request.URL.String()),
				slog.String("title", doc.DOM.Find("title").Text()),
			)
			return
		}

		extracted := false

		// selected html selections matching the query
		// might be empty if the query is not found
		for _, selected := range crawler.selections(doc) {

			ok, err := crawler.deps.Extractor.Extract(doc, selected)

			if err != nil {
				crawler.deps.Monitor.OnError(doc.Response, err)
				slog.Warn("extraction error",
					slog.String("error", err.Error()),
					slog.String("url", doc.Request.URL.String()),
					slog.String("title", doc.DOM.Find("title").Text()),
				)
				continue
			}

			if !ok {
				slog.Info("skipped",
					slog.String("url", doc.Request.URL.String()),
					slog.String("title", doc.DOM.Find("title").Text()),
				)
				continue
			}

			// send metrics
			crawler.deps.Monitor.OnExtract(doc.Response)
			crawler.CountExtraction()
			extracted = true
		}

		if !extracted {
			slog.Warn("no data extracted",
				slog.String("url", doc.Request.URL.String()),
				slog.String("title", doc.DOM.Find("title").Text()),
			)
		} else {
			slog.Info("extracted",
				slog.String("url", doc.Request.URL.String()),
				slog.String("title", doc.DOM.Find("title").Text()),
			)
		}
	}
}

// CountExtraction matching the query (with JS browse if Config.ExtractSelector is not found in GET response)
func (crawler *Dispatch) CountExtraction() {

	// add to the extracted count
	crawler.extractedCount.Add(1)

	// check if the limit is reached
	if crawler.IsExtractionLimitReached() {

		// stop the queue
		// existing requests will be processed
		// catch them on extraction with IsExtractionLimitReached
		crawler.queue.Stop()
	}
}

func (crawler *Dispatch) IsExtractionLimitReached() bool {
	count := int(crawler.extractedCount.Load())
	limit := crawler.args.ExtractLimit
	return limit > 0 && count >= limit
}

// selections matching the query (with JS browse if Config.ExtractSelector is not found in GET response)
func (crawler *Dispatch) selections(e *colly.HTMLElement) []*goquery.Selection {

	var browser Browser

	if crawler.args.UseBrowser {
		browser = crawler.browser
	}

	return Selections(e, crawler.args.ExtractSelector, browser)
}

// Selections matching the query (with JS browse if Config.ExtractSelector is not found in GET response)
func Selections(e *colly.HTMLElement, selector string, browser Browser) []*goquery.Selection {

	if selector == "html" {
		return []*goquery.Selection{e.DOM}
	}

	var selection *goquery.Selection

	if browser != nil {
		var err error
		if selection, err = browser.Browse(e.Request.URL.String()); err != nil {
			return nil
		}
	} else {
		selection = e.DOM.Find(selector)
	}

	var nodes []*goquery.Selection

	selection.Each(func(i int, s *goquery.Selection) {
		nodes = append(nodes, s)
	})

	return nodes
}

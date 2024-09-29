package events

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"net/url"
	"regexp"
)

// extract entries from html selections
func (crawler *Dispatch) extract() func(e *colly.HTMLElement) {

	var entityURL *regexp.Regexp

	// If the entity Endpoint is not empty, compile a regular expression from it
	if len(crawler.args.ExtractURL) > 0 {
		regex := config.RegexPattern(crawler.args.ExtractURL)
		entityURL = regexp.MustCompile(regex)
	}

	match := func(u *url.URL) bool {
		if entityURL == nil {
			return true
		}
		return entityURL.MatchString(u.String())
	}

	return func(doc *colly.HTMLElement) {

		// if expression exists, extract entity urls
		if !match(doc.Request.URL) {
			return
		}

		// selected html selections matching the query
		// might be empty if the query is not found
		for _, selected := range crawler.selections(doc) {

			if err := crawler.deps.Extractor(doc, selected); err != nil {
				crawler.deps.Monitor.OnError(doc.Response, err)
				continue
			}

			crawler.deps.Monitor.OnExtract(doc.Response)
			crawler.WatchLimit()
		}
	}
}

// WatchLimit matching the query (with JS browse if Config.ExtractSelector is not found in GET response)
func (crawler *Dispatch) WatchLimit() {

	count := int(crawler.extractedCount.Add(1))
	limit := crawler.args.ExtractLimit

	// limit is set and reached
	if limit > 0 && count >= limit {
		// stop the new requests
		crawler.queue.Stop()
	}
}

// selections matching the query (with JS browse if Config.ExtractSelector is not found in GET response)
func (crawler *Dispatch) selections(e *colly.HTMLElement) []*goquery.Selection {
	return Selections(e, crawler.args.ExtractSelector, crawler.browser)
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

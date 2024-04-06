package collect

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// collector based on colly
func (crawler *Crawler) collector() *colly.Collector {

	if crawler.collect != nil {
		return crawler.collect
	}

	// url regex from crawler args
	crawler.StartURL, crawler.AllowedURL, crawler.EntityURL = crawler.urlsRegexp()

	// default collector
	crawler.collect = colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(MustHost(crawler.StartURL), "127.0.0.1"),
		colly.UserAgent(crawler.UserAgent),
		colly.MaxDepth(crawler.Depth),
		colly.URLFilters(
			regexp.MustCompile(crawler.AllowedURL),
		),
	)

	// cookie handling
	// for turning off - crawler.collect.DisableCookies()
	j, err := cookiejar.New(&cookiejar.Options{})
	if err == nil {
		crawler.collect.SetCookieJar(j)
	}

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

	// entity url regex
	if len(crawler.EntityURL) > 0 {
		crawler._entityURL = regexp.MustCompile(crawler.EntityURL)
	}

	// set event handlers
	crawler.collect.OnHTML(`a[href]`, crawler.visit())
	crawler.collect.OnHTML(`html`, crawler.extract())

	return crawler.collect
}

// visit links found in the DOM
func (crawler *Crawler) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))

		if len(link) == 0 {
			return
		}

		if strings.HasPrefix(link, "#") {
			return
		}

		// skip images, scripts, etc.
		if !isValidURLExtension(link) {
			return
		}

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

		if errors.Is(err, colly.ErrForbiddenURL) {
			return
		}

		if errors.Is(err, colly.ErrNoURLFiltersMatch) {
			return
		}

		slog.Warn("visit failed",
			slog.String("url", link),
			slog.String("href", e.Attr("href")),
			slog.String("err", err.Error()),
		)
	}
}

// extract entries from html selections
func (crawler *Crawler) extract() func(e *colly.HTMLElement) {
	return func(doc *colly.HTMLElement) {

		// entity url regex
		if crawler._entityURL != nil {
			if !crawler._entityURL.MatchString(doc.Request.URL.String()) {
				return
			}
		}

		// selected html selections matching the query
		// might be empty if the query is not found
		for _, selected := range crawler.selections(doc) {
			err := crawler.Extractor(doc, selected)
			if err != nil {
				crawler.error(doc.Request.URL.String(), err)
				// explicitly
				continue
			}
		}
	}
}

func (crawler *Crawler) urlsRegexp() (start, allowed, entity string) {

	start = strings.TrimSpace(crawler.StartURL)
	if len(start) == 0 {
		panic("crawler: start url is required")
	}

	allowed = strings.TrimSpace(crawler.AllowedURL)
	if len(allowed) == 0 {
		// get the host from the start url
		allowed = MustRootUrl(start) + "{any}"
	}

	entityUrl := strings.TrimSpace(crawler.EntityURL)
	if len(entityUrl) == 0 {
		entityUrl = ""
	}

	return PlaceholdersToRegex(start), PlaceholdersToRegex(allowed), PlaceholdersToRegex(entityUrl)
}

func isValidURLExtension(urlStr string) bool {
	allowedExtensions := map[string]bool{
		".php":   true,
		".xhtml": true,
		".shtml": true,
		".cfm":   true,
		".html":  true,
		".htm":   true,
		".asp":   true,
		".aspx":  true,
		".jsp":   true,
		".jspx":  true,
	}

	// Parse the URL to extract the path
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Extract the file extension if present
	path := parsedURL.Path
	if dotIndex := strings.LastIndex(path, "."); dotIndex != -1 {
		ext := path[dotIndex:]
		return allowedExtensions[ext] // True if allowed, false otherwise
	}

	// True if no file extension is present
	return true
}

package collect

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/queue"
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

	// create a request queue with 2 consumer threads
	// https://go-colly.org/docs/examples/queue/
	var err error
	crawler.queue, err = queue.New(
		10, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 50000000}, // 50MB
	)

	// init metrics reporter
	crawler.report = NewReport()

	// url regex from crawler args
	crawler.StartURL, crawler.AllowedURL, crawler.EntityURL = crawler.urlsRegexp()

	// default collector
	crawler.collect = colly.NewCollector(
		colly.AllowedDomains(MustHost(crawler.StartURL), "127.0.0.1"),
		colly.MaxDepth(crawler.Depth),
		colly.URLFilters(regexp.MustCompile(crawler.AllowedURL)),
		colly.MaxBodySize(10<<20), // 10MB

		// colly.Async(true),
		// todo must be depending on crawl strategy chosen - singe or incremental
		// colly.AllowURLRevisit(),
	)

	// limit parallelism per domain
	//rule := &colly.LimitRule{DomainGlob: MustHost(crawler.StartURL), Parallelism: 1}
	//if err := crawler.collect.Limit(rule); err != nil {
	//	panic(err)
	//}

	extensions.RandomUserAgent(crawler.collect)

	// proxy handling
	if crawler.ProxyFn != nil {
		crawler.collect.SetProxyFunc(crawler.ProxyFn)
	}

	// timeouts
	crawler.collect.SetRequestTimeout(25 * time.Second)

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

	// entity url regex
	if len(crawler.EntityURL) > 0 {
		crawler._entityURL = regexp.MustCompile(crawler.EntityURL)
	}

	// Request setup
	crawler.retry = NewRetry()

	// set event handlers
	crawler.collect.OnHTML(`a[href]`, crawler.visit())
	crawler.collect.OnHTML(`html`, crawler.extract())
	crawler.collect.OnError(crawler.error)

	crawler.collect.OnRequest(func(r *colly.Request) {
		slog.Info("visiting", slog.String("url", r.URL.String()))
	})

	crawler.collect.OnResponse(func(r *colly.Response) {
		crawler.report.Visited()
	})

	return crawler.collect
}

// visit links found in the DOM
func (crawler *Crawler) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		// absolute url
		link := e.Request.AbsoluteURL(e.Attr("href"))

		// skip empty and anchor links
		if link == "" || strings.HasPrefix(link, "#") {
			return
		}

		// skip images, scripts, etc.
		if !isValidURLExtension(link) {
			return
		}

		// visit the link
		if err := crawler.queue.AddURL(link); err != nil {
			slog.Warn("crawler queue", slog.String("error", err.Error()))
		}
		//
		//err := crawler.collector().Visit(link)
		//if err == nil {
		//	crawler.report.Visited()
		//	return
		//}

		//// skip errors
		//skipErrors := []error{
		//	colly.ErrAlreadyVisited,
		//	colly.ErrForbiddenDomain,
		//	colly.ErrForbiddenURL,
		//	colly.ErrNoURLFiltersMatch,
		//}
		//for _, skip := range skipErrors {
		//	if errors.Is(err, skip) {
		//		slog.Debug("ignore error", slog.String("error", err.Error()))
		//		return
		//	}
		//}
		//
		//// log the error
		//slog.Warn("crawler visit",
		//	slog.String("url", link),
		//	slog.String("proxy", e.Request.ProxyURL),
		//	slog.String("error", err.Error()),
		//)
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
				crawler.report.ExtractFailed()
				continue
			}

			crawler.report.Extracted()
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
		allowed := allowedExtensions[ext] // True if allowed, false otherwise
		return allowed
	}

	// True if no file extension is present
	return true
}
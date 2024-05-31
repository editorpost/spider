package collect

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/queue"
	"github.com/gocolly/colly/v2/storage"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// setup based on colly
func (crawler *Crawler) setup() *colly.Collector {

	// return if already initialized
	if crawler.collect != nil {
		return crawler.collect
	}

	// init metrics reporter
	crawler.report = NewReport()
	crawler.monitor = NewMetrics(crawler.JobID, crawler.SpiderID)

	// create a request queue with 2 consumer threads
	// https://go-colly.org/docs/examples/queue/
	crawler.withQueue()
	crawler.withCollector()
	crawler.withProxy()
	crawler.withEventHandlers()

	return crawler.collect
}

// withCollector initializes the collector for the crawler.
// It first checks if the collector is already initialized, if so, it returns the existing collector.
// If not, it sets up a new collector with a maximum depth and maximum body size.
// It then applies URL filters and storage to the collector and sets a random user agent.
// Finally, it returns the initialized collector.
//
// Returns:
//
//	*colly.Collector: The initialized collector.
func (crawler *Crawler) withCollector() *colly.Collector {

	// Check if the collector is already initialized
	if crawler.collect != nil {
		return crawler.collect
	}

	// Set up a new collector with a maximum depth and maximum body size
	crawler.collect = colly.NewCollector(
		colly.MaxDepth(crawler.Depth),
		colly.MaxBodySize(10<<20), // 10MB
	)

	// Apply URL filters and storage to the collector
	crawler.withURLFilters()
	crawler.withStorage()

	// Set a random user agent
	extensions.RandomUserAgent(crawler.collect)

	// Return the initialized collector
	return crawler.collect
}

// withStorage sets up the storage for the crawler
// or creates an in-memory storage if not provided.
func (crawler *Crawler) withStorage() {

	// Check if the Storage field of the Crawler struct is not nil
	if crawler.Storage == nil {
		crawler.Storage = &storage.InMemoryStorage{}
		return
	}

	// Retry to set the Storage as the storage for the collector
	err := crawler.collect.SetStorage(crawler.Storage)
	// If an error occurs, panic and stop the execution
	if err != nil {
		panic(err)
	}
}

// withURLFilters sets up the URL filters for the crawler.
// It first generates regular expressions from the start, allowed, and entity URLs.
// If the entity URL is not empty, it compiles a regular expression from it.
// It then appends the host of the start URL to the allowed domains of the collector.
// Finally, it appends a regular expression from the allowed URL to the URL filters of the collector.
//
// This function does not return a value.
func (crawler *Crawler) withURLFilters() {

	// Generate regular expressions from the start, allowed, and entity URLs
	crawler.StartURL, crawler.AllowedURL, crawler.EntityURL = crawler.urlsRegexp()

	// If the entity URL is not empty, compile a regular expression from it
	if len(crawler.EntityURL) > 0 {
		crawler._entityURL = regexp.MustCompile(crawler.EntityURL)
	}

	// Append the host of the start URL to the allowed domains of the collector
	crawler.collect.AllowedDomains = append(crawler.collect.AllowedDomains, MustHostname(crawler.StartURL))

	// Append a regular expression from the allowed URL to the URL filters of the collector
	crawler.collect.URLFilters = append(crawler.collect.URLFilters, regexp.MustCompile(crawler.AllowedURL))
}

// withQueue sets up the request queue for the crawler.
// It creates a new request queue with 25 consumer threads and an in-memory queue storage with a maximum size of 50MB.
// If an error occurs during the setup, it panics and stops the execution.
//
// This function does not return a value.
func (crawler *Crawler) withQueue() {
	// create a request queue with 25 consumer threads
	// https://go-colly.org/docs/examples/queue/
	var err error
	crawler.queue, err = queue.New(
		5, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 5000000}, // 5MB
	)
	if err != nil {
		panic(err)
	}
}

// inMemoryQueueItem hold urls max len 512 bytes
// e.g. 5MB can hold 10k urls
type inMemoryQueueItem struct {
	Request []byte
	Next    *inMemoryQueueItem
}

// withProxy sets up the proxy for the crawler.
// It sets the request timeout for the collector to 25 seconds.
// It then sets up a cookie jar for the collector, if an error occurs during the setup, it skips this step.
// Finally, it sets up retries for response and proxy errors.
func (crawler *Crawler) withProxy() {

	// round tripper
	if crawler.RoundTripper != nil {
		crawler.collect.WithTransport(crawler.RoundTripper)
	}

	// proxy func, note this injects to transport
	// it is better to call after transport init.
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

	// retries
	crawler.errRetry = NewRetry(ResponseRetries)
	crawler.proxyRetry = NewRetry(BadProxyRetries)
}

func (crawler *Crawler) withDebug() {

	if crawler.Debugger == nil {
		return
	}

	// colly event dispatcher for logging, monitoring and debugging
	crawler.collect.SetDebugger(crawler.Debugger)
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

package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
)

type Config struct {

	// StartURL is the url to start the scraping
	StartURL string `json:"StartURL"`

	// AllowedURLs is comma separated regex to match the urls
	// use it to reduce the number of urls to visit
	AllowedURLs []string `json:"AllowedURLs"`

	// ExtractURLs is the regex to match the entity urls
	// use it to extract the entity urls
	ExtractURLs []string `json:"ExtractURLs"`

	// VisitOnce is the flag to visit the url only once
	// If true, then visited urls stored in S3 between runs
	// Else default colly memory collector used to store the visited urls by default
	// def: false
	VisitOnce bool `json:"VisitOnce"`

	// ExtractSelector is the css selector to match the elements
	// use selector for extracting entities and filtering pages
	// def: html
	ExtractSelector string `json:"ExtractSelector"`

	// ExtractLimit is the limit of entities to extract
	// Crawler gracefully stops after reaching the limit
	ExtractLimit int `json:"ExtractLimit"`

	// UseBrowser is a flag to use browser for rendering the page
	UseBrowser bool `json:"UseBrowser"`

	// Depth if is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	Depth int `json:"Depth"`

	// UserAgent is the user agent string used by the collector
	UserAgent string `json:"UserAgent"`

	// ProxyEnabled is the flag to enable proxy or send requests directly
	ProxyEnabled bool `json:"ProxyEnabled"`

	// ProxySources is the list of proxy sources, expected to return list of proxies URLs.
	// If empty, the default proxy sources is used.
	ProxySources []string `json:"ProxySources"`
}

// The Config JSON representation:
// {
// 	"ID": "ready-check",
//  "Name": "Ready Check",
// 	"StartURL": "https://example.com",
// 	"AllowedURL": "https://example.com/{any}",
// 	"ExtractURL": "https://example.com/articles/{any}",
// 	"ExtractSelector": "article",
// 	"ExtractLimit": 1,
// 	"UseBrowser": true,
// 	"Depth": 1,
// 	"UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
// 	"ProxyEnabled": true,
// 	"ProxySources": []
// }

func (args *Config) Normalize() error {

	if err := args.NormalizeURLs(); err != nil {
		return err
	}

	if err := args.NormalizeUserAgent(); err != nil {
		return err
	}

	args.NormalizeExtractSelector()

	return nil
}

func (args *Config) Log() slog.Attr {
	return slog.Group("args",
		slog.String("start_url", args.StartURL),
		slog.String("allowed_urls", strings.Join(args.AllowedURLs, ",")),
		slog.String("extract_urls", strings.Join(args.ExtractURLs, ",")),
		slog.String("entity_selector", args.ExtractSelector),
		slog.Bool("use_browser", args.UseBrowser),
		slog.Int("depth", args.Depth),
		slog.String("user_agent", args.UserAgent),
	)
}

func (args *Config) NormalizeExtractSelector() {
	args.ExtractSelector = strings.TrimSpace(args.ExtractSelector)
	if len(args.ExtractSelector) == 0 {
		args.ExtractSelector = "html"
	}
}

func (args *Config) NormalizeURLs() error {

	// start url is required
	args.StartURL = strings.TrimSpace(args.StartURL)
	if len(args.StartURL) == 0 {
		return errors.New("start url is required")
	}

	// start url should be valid
	startURI, err := url.ParseRequestURI(args.StartURL)
	if err != nil {
		return fmt.Errorf("start url is invalid: %w", err)
	}

	// if host is empty, then it is invalid
	if len(startURI.Host) == 0 {
		return errors.New("start url host is invalid, add domain name")
	}

	// by default all urls are allowed including main page
	if len(args.AllowedURLs) == 0 {
		// no slash separator between root url and any
		// to include main page w/o trailing slash.
		args.AllowedURLs = []string{RootUrl(startURI) + "{any}"}
	}

	return nil
}

// NormalizeUserAgent sets the default user agent
func (args *Config) NormalizeUserAgent() error {

	args.UserAgent = strings.TrimSpace(args.UserAgent)
	if len(args.UserAgent) == 0 {
		args.UserAgent = DefaultUserAgent
	}

	return nil
}

// RootUrl return the root url
// e.g. https://example.com/articles/1234/5678 => https://example.com
func RootUrl(u *url.URL) string {
	return u.Scheme + "://" + u.Host // no port explicitly
}

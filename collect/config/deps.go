package config

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/storage"
	"net/http"
)

type Deps struct {
	// RoundTripper is the function to return the next proxy from the list
	RoundTripper http.RoundTripper
	// Extractor is the function to process matched the data, e.g. html tag node
	Extractor func(*colly.HTMLElement, *goquery.Selection) error
	// Storage is the storage backend for the collector
	Storage storage.Storage
	// Debugger is the function to debug/Monitor spider
	Debugger debug.Debugger
	// Metrics is the spider event dispatcher and VictoriaMetrics
	Monitor Metrics
}

// Normalize default values
func (deps *Deps) Normalize() *Deps {

	if deps.Extractor == nil {
		deps.Extractor = func(e *colly.HTMLElement, s *goquery.Selection) error {
			return nil
		}
	}

	if deps.Monitor == nil {
		deps.Monitor = &MetricsFallback{}
	}

	if deps.Storage == nil {
		deps.Storage = &storage.InMemoryStorage{}
	}

	return deps
}

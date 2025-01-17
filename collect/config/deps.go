package config

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/storage"
	"net/http"
)

type Extractor interface {
	Extract(*colly.HTMLElement, *goquery.Selection) (bool, error)
}

type Deps struct {
	// RoundTripper is the function to return the next proxy from the list
	RoundTripper http.RoundTripper
	// Extractor is the function to process matched the data, e.g. html tag node
	Extractor Extractor
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
		deps.Extractor = &ExtractorFn{}
	}

	if deps.Monitor == nil {
		deps.Monitor = &MetricsFallback{}
	}

	if deps.Storage == nil {
		deps.Storage = &storage.InMemoryStorage{}
	}

	return deps
}

type ExtractorFn struct {
	fn func(_ *colly.HTMLElement, _ *goquery.Selection) (bool, error)
}

func NewExtractor(fn func(_ *colly.HTMLElement, _ *goquery.Selection) (bool, error)) *ExtractorFn {
	return &ExtractorFn{fn: fn}
}

func (e *ExtractorFn) Extract(doc *colly.HTMLElement, selection *goquery.Selection) (bool, error) {
	if e.fn == nil {
		return false, nil
	}
	return e.fn(doc, selection)
}

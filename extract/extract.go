package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net/url"
	"strings"
)

const (
	UrlField  = "spider__url"
	HostField = "spider__host"
	HtmlField = "spider__html"
)

//goland:noinspection GoNameStartsWithPackageName
type (
	PipeFn    func(*Payload) error
	ExtractFn func(doc *colly.HTMLElement, s *goquery.Selection) error

	Payload struct {
		// Doc is full document
		Doc *colly.HTMLElement `json:"-"`
		// Selection of entity in document
		Selection *goquery.Selection `json:"-"`
		// URL of the document
		URL *url.URL `json:"-"`
		// Data is a map of extracted data
		Data map[string]any `json:"Data"`
	}
)

var (
	// ErrDataNotFound expected error, stops the extraction pipeline.
	ErrDataNotFound = errors.New("skip entity extraction, required data is missing")
)

// Pipe is a function to process the html node and url.
// Order of the pipes is important.
func Pipe(pipes ...PipeFn) ExtractFn {

	return func(doc *colly.HTMLElement, s *goquery.Selection) error {

		if s == nil {
			return errors.New("document is nil")
		}

		payload := &Payload{
			Doc:       doc,
			Selection: s,
			URL:       doc.Request.URL,
			Data: map[string]any{
				UrlField:  doc.Request.URL.String(),
				HostField: doc.Request.URL.Host,
			},
		}

		for _, pipe := range pipes {
			err := pipe(payload)

			// stop the pipe chain if required data is missing
			if errors.Is(err, ErrDataNotFound) {
				return nil
			}

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Extractors(ff []*fields.Field, entities string) ([]PipeFn, error) {

	// entity extractors
	extractors := ExtractorsByName(entities)

	// field extractors
	extractFields, err := Fields(ff...)
	if err != nil {
		slog.Error("build extractors from field tree", err)
		return nil, err
	}

	extractors = append(extractors, extractFields)

	return extractors, err
}

// ExtractorsByName creates slice of extractors by name.
// The string is a string like "html,article", e.g.: extract.Html, extract.Article
func ExtractorsByName(seq string) []PipeFn {

	seq = strings.ReplaceAll(seq, " ", "")

	if seq == "" {
		return []PipeFn{}
	}

	extractors := make([]PipeFn, 0)
	for _, key := range strings.Split(seq, ",") {

		switch key {
		case "html":
			extractors = append(extractors, Html)
		case "article":
			extractors = append(extractors, Article)
		}
	}

	return extractors
}

// ExtractorsByJsonString creates slice of extractors by name.
//
//goland:noinspection GoUnusedExportedFunction
func ExtractorsByJsonString(js string) []PipeFn {
	if js == "" {
		return []PipeFn{}
	}
	return []PipeFn{}
}

package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strings"
)

const (
	UrlField  = "spider__url"
	HostField = "spider__host"
	HtmlField = "spider__html"
)

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
			if err != nil {
				return err
			}
		}

		return nil
	}
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
func ExtractorsByJsonString(js string) []PipeFn {
	if js == "" {
		return []PipeFn{}
	}
	return []PipeFn{}
}

package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
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
		Doc *colly.HTMLElement
		// Selection of entity in document
		Selection *goquery.Selection
		// URL of the document
		URL *url.URL
		// Data is a map of extracted data
		Data map[string]any
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

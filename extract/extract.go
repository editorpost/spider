package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
)

const (
	EntityArticle = "article"
	UrlField      = "url"
	HtmlField     = "html"
)

type (
	PipeFn    func(*Payload) error
	ExtractFn func(doc *colly.HTMLElement, s *goquery.Selection) error

	Payload struct {
		Doc       *colly.HTMLElement
		Selection *goquery.Selection
		URL       *url.URL
		Data      map[string]any
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
				UrlField: doc.Request.URL.String(),
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

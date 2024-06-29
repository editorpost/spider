package payload

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/media"
	"github.com/gocolly/colly/v2"
	"net/url"
)

const (
	HtmlField = "spider__html"
	UrlField  = "spider__url"
	HostField = "spider__host"
)

var (
	// ErrDataNotFound expected error, stops the extraction pipeline.
	ErrDataNotFound = errors.New("skip entity extraction, required data is missing")
)

type (
	Extractor func(*Payload) error
	//goland:noinspection GoNameStartsWithPackageName
	Payload struct {
		// Doc is full document
		Doc *colly.HTMLElement `json:"-"`
		// Selection of entity in document
		Selection *goquery.Selection `json:"-"`
		// URL of the document
		URL *url.URL `json:"-"`
		// Data is a map of extracted data
		Data map[string]any `json:"Data"`
		// Claims is a list of media to upload
		Claims *media.Claims `json:"Claims"`
	}
	CollectorHook func(doc *colly.HTMLElement, s *goquery.Selection) error
)

type Pipeline struct {
	extractors []Extractor
}

// PipelineFn of Processors. Order matters.
func PipelineFn(extractors ...Extractor) CollectorHook {

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

		for _, extractor := range extractors {
			err := extractor(payload)

			// stop the extractor chain if required data is missing
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

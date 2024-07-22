package pipe

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"hash/fnv"
	"net/url"
	"time"
)

const (
	SpiderIDField = "spider__id"
	HtmlField     = "spider__html"
	UrlField      = "spider__url"
	HostField     = "spider__host"
	DateField     = "spider__date"
)

var (
	// ErrDataNotFound expected error, stops the extraction pipeline.
	ErrDataNotFound = errors.New("skip entity extraction, required data is missing")
)

type (
	Extractor func(*Payload) error
	//goland:noinspection GoNameStartsWithPackageName
	Payload struct {
		// ID is document url hash
		ID  string          `json:"ID"`
		Ctx context.Context `json:"-"`
		// Doc is full document
		Doc *colly.HTMLElement `json:"-"`
		// Selection of entity in document
		Selection *goquery.Selection `json:"-"`
		// URL of the document
		URL *url.URL `json:"-"`
		// Data is a map of extracted data
		Data map[string]any `json:"Data"`
	}
	CollectorHook func(doc *colly.HTMLElement, s *goquery.Selection) error
)

func NewPayload(doc *colly.HTMLElement, s *goquery.Selection) (*Payload, error) {

	if s == nil {
		return nil, errors.New("document is nil")
	}

	id, err := Hash(doc.Request.URL.String())
	if err != nil {
		return nil, fmt.Errorf("url FNV hash error: %w", err)
	}

	return &Payload{
		ID:        id,
		Ctx:       context.Background(),
		Doc:       doc,
		Selection: s,
		URL:       doc.Request.URL,
		Data: map[string]any{
			SpiderIDField: id,
			DateField:     time.Now().UTC().String(),
			HostField:     doc.Request.URL.Host,
			UrlField:      doc.Request.URL.String(),
		},
	}, nil
}

// Hash generates an FNV hash from the source Endpoint.
func Hash(uri string) (string, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum64()), nil
}

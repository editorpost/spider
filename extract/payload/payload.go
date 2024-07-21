package payload

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
		ID  string
		Ctx context.Context
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

type Pipeline struct {
	// starter extractors called before the main extractors
	starter []Extractor
	// finisher extractors called after the main extractors
	finisher []Extractor
	// extractors is a list of main extractors
	extractors []Extractor
}

func NewPipeline(extractors ...Extractor) *Pipeline {
	return &Pipeline{
		extractors: extractors,
		starter:    make([]Extractor, 0),
		finisher:   make([]Extractor, 0),
	}
}

func (p *Pipeline) Append(extractors ...Extractor) *Pipeline {
	p.extractors = append(p.extractors, extractors...)
	return p
}

func (p *Pipeline) Starter(extractors ...Extractor) *Pipeline {
	p.starter = append(p.starter, extractors...)
	return p
}

func (p *Pipeline) Finisher(extractors ...Extractor) *Pipeline {
	p.finisher = append(p.finisher, extractors...)
	return p
}

func (p *Pipeline) Extract(doc *colly.HTMLElement, s *goquery.Selection) error {

	if s == nil {
		return errors.New("document is nil")
	}

	id, err := Hash(doc.Request.URL.String())
	if err != nil {
		return fmt.Errorf("url FNV hash error: %w", err)
	}

	payload := &Payload{
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
	}

	// starter
	if err := p.exec(payload, p.starter...); err != nil {
		return err
	}

	// main
	if err := p.exec(payload, p.extractors...); err != nil {
		return err
	}

	// finisher
	if err := p.exec(payload, p.finisher...); err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) exec(payload *Payload, extractors ...Extractor) error {

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

// PipelineFn of Processors. Order matters.
// @deprecated
func PipelineFn(extractors ...Extractor) CollectorHook {

	return func(doc *colly.HTMLElement, s *goquery.Selection) error {

		if s == nil {
			return errors.New("document is nil")
		}

		id, err := Hash(doc.Request.URL.String())
		if err != nil {
			return fmt.Errorf("url FNV hash error: %w", err)
		}

		payload := &Payload{
			ID:        id,
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

// Hash generates an FNV hash from the source Endpoint.
func Hash(uri string) (string, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum64()), nil
}

package pipe

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
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

	payload, err := NewPayload(doc, s)
	if err != nil {
		return fmt.Errorf("payload creation error: %w", err)
	}

	// set job id and provider
	if err = JobMetadata(payload); err != nil {
		return err
	}

	// starter
	if err = p.exec(payload, p.starter...); err != nil {
		return err
	}

	// main
	if err = p.exec(payload, p.extractors...); err != nil {
		return err
	}

	// finisher
	if err = p.exec(payload, p.finisher...); err != nil {
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

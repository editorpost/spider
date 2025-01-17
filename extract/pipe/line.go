package pipe

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type (
	Pipeline struct {
		// starter extractors called before the main extractors
		starter []Extractor
		// finisher extractors called after the main extractors
		finisher []Extractor
		// extractors is a list of main extractors
		extractors []Extractor
		// history keeps extraction history, avoiding/allowing duplication
		history *History
	}
)

func NewPipeline(extractors ...Extractor) *Pipeline {
	return &Pipeline{
		extractors: extractors,
		starter:    make([]Extractor, 0),
		finisher:   make([]Extractor, 0),
		history:    NewPayloadHistory(),
	}
}

func (p *Pipeline) History() *History {
	return p.history
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

func (p *Pipeline) Extract(doc *colly.HTMLElement, s *goquery.Selection) (extracted bool, err error) {

	payload, err := NewPayload(doc, s)
	if err != nil {
		return false, fmt.Errorf("payload creation error: %w", err)
	}

	// check if payload is already extracted
	if extracted, err = p.history.IsExtracted(payload); err != nil || extracted {
		return
	}

	// set job id and provider
	if err = JobMetadata(payload); err != nil {
		return
	}

	// starter
	if err = p.exec(payload, p.starter...); err != nil {
		return
	}

	// main
	if err = p.exec(payload, p.extractors...); err != nil {
		return
	}

	// finisher
	if err = p.exec(payload, p.finisher...); err != nil {
		return
	}

	return true, nil
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

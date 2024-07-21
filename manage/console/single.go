package console

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage/setup"
)

// Single spider against one url and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Single(uri, selector string, extractor payload.Extractor) (*payload.Payload, error) {

	result := &payload.Payload{}

	args := &config.Args{
		// Any name since no data is stored
		ID: "ready-check",
		// All urls are the same for single turn
		StartURL:        uri,
		AllowedURL:      uri,
		ExtractURL:      uri,
		ExtractSelector: selector,
		ExtractLimit:    1,
		// Depth is 1 for single turn
		Depth:        1,
		ProxyEnabled: true,
	}

	s, err := setup.NewSpider(args, &extract.Config{})
	if err != nil {
		return result, err
	}

	s.Pipeline().Append(extractor)
	s.Pipeline().Finisher(func(p *payload.Payload) error {
		result = p
		return nil
	})

	// empty deploy, no data stored
	crawler, err := s.NewCrawler(&setup.Deploy{})
	if err != nil {
		return result, err
	}

	defer s.Shutdown()
	return result, crawler.Run()
}

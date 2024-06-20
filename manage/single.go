package manage

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/setup"
)

// Single spider against one url and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Single(uri, selector string, extractor extract.PipeFn) (*extract.Payload, error) {

	result := &extract.Payload{}

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

	// empty deploy, since no data is stored
	deploy := &setup.Config{}

	err := Start(args, deploy, extract.Article, extractor)

	return result, err
}

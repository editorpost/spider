package manage

import (
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/setup"
)

// SingleTurn spider takes a single URL and extracts the data
// It does not store the data, but uses proxy pool for requests.
func SingleTurn(uri, selector string, extractor extract.PipeFn) (*extract.Payload, error) {

	result := &extract.Payload{}

	args := &setup.Config{
		// Any name since no data is stored
		Name: "ready-check",
		// All urls are the same for single turn
		StartURL:   uri,
		AllowedURL: uri,
		EntityURL:  uri,
		// Depth is 1 for single turn
		Depth:          1,
		EntitySelector: selector,
		// EntityExtract executes the end of the pipeline due there is no storage.
		// Copying payload on the final step to the result variable.
		EntityExtract: func(payload *extract.Payload) error {
			result = payload
			return extractor(payload)
		},
		// Used default in-memory storage
	}

	err := Start(args)
	return result, err
}

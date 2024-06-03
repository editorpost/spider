package windmill

import (
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(argsMap any, extractor extract.PipeFn) ([]*extract.Payload, error) {

	args := &config.Args{}
	err := script.ParseArgs(argsMap, args)
	if err != nil {
		return nil, err
	}

	args.SpiderID = "trial"
	items := []*extract.Payload{}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	limiter := func(payload *extract.Payload) error {
		if len(items) < args.ExtractLimit {
			items = append(items, payload)
		}

		return nil
	}

	err = Start(argsMap, extractor, limiter)

	return items, err
}

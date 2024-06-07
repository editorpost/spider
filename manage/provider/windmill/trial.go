package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(argsMap any, extractors ...extract.PipeFn) ([]*extract.Payload, error) {

	args := &config.Args{
		SpiderID: "trials",
	}

	err := vars.FromJSON(argsMap, args)
	if err != nil {
		return nil, err
	}

	items := []*extract.Payload{}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	limiter := func(payload *extract.Payload) error {
		if len(items) < args.ExtractLimit {
			items = append(items, payload)
		}

		return nil
	}

	err = Start(argsMap, append(extractors, limiter)...)

	return items, err
}

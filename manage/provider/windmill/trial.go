package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(args *config.Args, extractors ...extract.PipeFn) error {

	items := []*extract.Payload{}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	limiter := func(payload *extract.Payload) error {
		if len(items) < args.ExtractLimit {
			items = append(items, payload)
		}
		return nil
	}

	if err := Start(args, append(extractors, limiter)...); err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects

	return vars.WriteScriptResult(items, "./result.json")
}

//goland:noinspection GoUnusedExportedFunction
func TrialWith(argsMap any, extractors ...extract.PipeFn) error {

	args := &config.Args{
		SpiderID: "trials",
	}

	if err := vars.FromJSON(argsMap, args); err != nil {
		return err
	}

	return Trial(args, extractors...)
}

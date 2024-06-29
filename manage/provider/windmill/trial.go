package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(args *config.Args, extractors ...payload.Extractor) error {

	args.ID = "trial"
	items := []map[string]any{}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	limiter := func(payload *payload.Payload) error {
		if len(items) < args.ExtractLimit {
			items = append(items, payload.Data)
		}
		return nil
	}

	if err := manage.Start(args, &setup.Config{}, append(extractors, limiter)...); err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects
	return vars.WriteScriptResult(items, "./result.json")
}

//goland:noinspection GoUnusedExportedFunction
func TrialWith(argsMap any, extractors ...payload.Extractor) error {

	args := &config.Args{
		ID: "trials",
	}

	if err := vars.FromJSON(argsMap, args); err != nil {
		return err
	}

	return Trial(args, extractors...)
}

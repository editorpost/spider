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
func Trial(args *config.Args, pipe *payload.Pipeline) error {

	args.ID = "trial"
	items := []map[string]any{}

	// force low limit for trial
	if args.ExtractLimit == 0 && args.ExtractLimit > 30 {
		args.ExtractLimit = 30
	}

	pipe.Append(func(payload *payload.Payload) error {
		items = append(items, payload.Data)
		return nil
	})

	if err := manage.Start(args, &setup.Config{}, pipe); err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects
	return vars.WriteScriptResult(items, "./result.json")
}

//goland:noinspection GoUnusedExportedFunction
func TrialWith(argsMap any, pipe *payload.Pipeline) error {

	args := &config.Args{
		ID: "trials",
	}

	if err := vars.FromJSON(argsMap, args); err != nil {
		return err
	}

	return Trial(args, pipe)
}

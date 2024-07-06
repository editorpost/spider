package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(s *setup.Spider) error {

	s.Args.ID = "trial"
	items := []map[string]any{}

	// force low limit for trial
	if s.Args.ExtractLimit == 0 && s.Args.ExtractLimit > 30 {
		s.Args.ExtractLimit = 30
	}

	s.Pipeline().Append(func(payload *payload.Payload) error {
		items = append(items, payload.Data)
		return nil
	})

	if err := manage.Start(s, &setup.Deploy{}); err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects
	return vars.WriteScriptResult(items, "./result.json")
}

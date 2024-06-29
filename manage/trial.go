package manage

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(args *config.Args, pipe *payload.Pipeline, extractor payload.Extractor) ([]*payload.Payload, error) {

	args.ID = "trial"
	items := []*payload.Payload{}
	deploy := &setup.Config{}

	// force low limit for trial
	if args.ExtractLimit == 0 && args.ExtractLimit > 30 {
		args.ExtractLimit = 30
	}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	pipe.Append(func(payload *payload.Payload) error {
		items = append(items, payload)
		return nil
	})

	err := Start(args, deploy, pipe)

	return items, err
}

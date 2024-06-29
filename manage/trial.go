package manage

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(args *config.Args, extractor payload.Extractor) ([]*payload.Payload, error) {

	args.ID = "trial"
	items := []*payload.Payload{}
	deploy := &setup.Config{}

	err := Start(args, deploy, extractor, func(payload *payload.Payload) error {

		// the queue will stop automatically
		// after args.ExtractLimit is reached
		if len(items) < args.ExtractLimit {
			items = append(items, payload)
		}

		return nil
	})

	return items, err
}

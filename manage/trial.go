package manage

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(args *config.Args, extractor extract.PipeFn) ([]*extract.Payload, error) {

	args.SpiderID = "trial"
	items := []*extract.Payload{}
	deploy := &setup.Deploy{}

	err := Start(args, deploy, extractor, func(payload *extract.Payload) error {

		// the queue will stop automatically
		// after args.ExtractLimit is reached
		if len(items) < args.ExtractLimit {
			items = append(items, payload)
		}

		return nil
	})

	return items, err
}

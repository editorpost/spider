package manage

import (
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage/setup"
)

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(s *setup.Spider) ([]*payload.Payload, error) {

	s.Args.ID = "trial"
	items := []*payload.Payload{}

	// force low limit for trial
	if s.Args.ExtractLimit == 0 && s.Args.ExtractLimit > 30 {
		s.Args.ExtractLimit = 30
	}

	// the queue will stop automatically
	// after args.ExtractLimit is reached
	s.Pipeline().Append(func(payload *payload.Payload) error {
		items = append(items, payload)
		return nil
	})

	err := Start(s, &setup.Config{})

	return items, err
}

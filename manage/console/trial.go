package console

import (
	"errors"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/manage/setup"
)

// Trial runs spider and return extracted data without storing it.
// It is allowed to use proxy pool for requests.
func Trial(s *setup.Spider) ([]*pipe.Payload, error) {

	// Spider.ID is required
	if s.ID == "" {
		return nil, errors.New("spider check: ID not provided")
	}

	items := []*pipe.Payload{}

	// force low limit for trial
	if s.Collect.ExtractLimit == 0 || s.Collect.ExtractLimit > 30 {
		s.Collect.ExtractLimit = 30
	}

	s.Pipeline().Append(func(payload *pipe.Payload) error {
		items = append(items, payload)
		return nil
	})

	crawler, err := s.NewCrawler()
	if err != nil {
		return items, err
	}

	defer s.Shutdown()
	if err = crawler.Run(); err != nil {
		return items, err
	}

	return items, err
}

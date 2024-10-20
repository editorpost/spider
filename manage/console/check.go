package console

import (
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
)

// Check runs spider as usual, but with limited extract limit and storage paths.
// It is used for testing spider configuration and extractors from the console, api or clients.
func Check(checkID string, s *setup.Spider) (string, error) {

	// Replace actual spider ID with Windmill job UUID,
	// making virtual copy of the actual spider configuration,
	// keeping relation between Windmill job and check run.
	// It guarantees that check runs are isolated,
	// but close to real runs.
	s.ID = checkID

	// force low hard-limit for check runs
	if s.Collect.ExtractLimit == 0 || s.Collect.ExtractLimit > 30 {
		s.Collect.ExtractLimit = 30
	}

	// replace actual storage paths with check storage paths
	s.Deploy.Paths = store.CheckStoragePaths()

	// return check UUID and run the spider
	return s.ID, Start(s)
}

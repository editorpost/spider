package console

import (
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
)

// Check runs spider as usual, but with limited extract limit and storage paths.
// It is used for testing spider configuration and extractors from the console, api or clients.
func Check(spider *setup.Spider) (map[string]any, error) {

	// force low hard-limit for check runs
	if spider.Collect.ExtractLimit == 0 || spider.Collect.ExtractLimit > 30 {
		spider.Collect.ExtractLimit = 30
	}

	// replace actual storage paths with check storage paths
	spider.Deploy.Paths = store.CheckStoragePaths()

	// return check UUID and run the spider
	return map[string]any{
		"CheckID": spider.ID,
		"Paths":   spider.Deploy.Paths,
	}, Start(spider)
}

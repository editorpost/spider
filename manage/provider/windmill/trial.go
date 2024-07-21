package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/manage/setup"
)

const JobResultFile = "./result.json"

// Trial spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Trial(s *setup.Spider) error {

	items, err := console.Trial(s)
	if err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects
	return vars.WriteScriptResult(items, JobResultFile)
}

package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/manage/setup"
)

const JobResultFile = "./result.json"

// Check spider against limited and return extracted data
// It does not store the data, but uses proxy pool for requests.
func Check(spider *setup.Spider) error {

	// result contains check ID and storage paths
	result, err := console.Check(spider)
	if err != nil {
		return err
	}

	// write extracted data to `./result.json` as windmill expects
	return vars.WriteScriptResult(result, JobResultFile)
}

package windmill

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/manage/setup"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
//
//goland:noinspection GoDfaNilDereference,GoUnusedExportedFunction
func Reset(spiderID string) error {

	var deploy *setup.Deploy

	if err := LoadDeployResource(deploy); err != nil {
		return err
	}

	return console.Reset(spiderID, deploy)
}

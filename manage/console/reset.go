package console

import (
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
//
//goland:noinspection GoDfaNilDereference,GoUnusedExportedFunction
func Reset(spiderID string, deploy *setup.Deploy) error {

	// ResetMedia drops the media store
	// All media files will be erased.
	if err := store.DropMediaStorage(spiderID, deploy.Media.S3); err != nil {
		return err
	}

	// ResetCollector drops the collector store
	// Crawler Endpoint history will be erased.
	if err := store.DropCollectStorage(spiderID, deploy.Storage); err != nil {
		return err
	}

	// ResetExtractor drops the extractor store
	// Extracted data will be erased. All temporary data/images will be lost.
	return store.DropExtractStorage(spiderID, deploy.Storage)
}

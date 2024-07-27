package console

import (
	"github.com/editorpost/donq/res"
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
//
//goland:noinspection GoDfaNilDereference,GoUnusedExportedFunction
func Reset(spiderID string, deploy *setup.Deploy) error {

	if err := ResetMedia(spiderID, deploy.Media); err != nil {
		return err
	}

	if err := ResetCollector(spiderID, deploy.Storage); err != nil {
		return err
	}

	return ResetExtractor(spiderID, deploy.Storage)
}

// ResetCollector drops the collector store
// Crawler Endpoint history will be erased.
func ResetCollector(spiderID string, bucket *res.S3) error {

	collector, _, err := store.NewCollectStorage(spiderID, bucket)
	if err != nil {
		return err
	}

	return collector.Reset()
}

// ResetExtractor drops the extractor store
// Extracted data will be erased. All temporary data/images will be lost.
func ResetExtractor(spiderID string, bucket *res.S3) error {

	extractor, err := store.NewExtractStorage(spiderID, bucket)
	if err != nil {
		return err
	}

	return extractor.Reset()
}

// ResetMedia drops the media store
// All media files will be erased.
func ResetMedia(spiderID string, bucket *res.S3Public) error {

	media, err := store.NewMediaStorage(spiderID, &bucket.S3)
	if err != nil {
		return err
	}

	return media.Reset()
}

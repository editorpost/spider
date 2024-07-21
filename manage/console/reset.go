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

	if err := ResetMedia(spiderID, deploy.Bucket); err != nil {
		return err
	}

	if err := ResetCollector(spiderID, deploy.Bucket); err != nil {
		return err
	}

	return ResetExtractor(spiderID, deploy.Bucket)
}

// ResetCollector drops the collector store
// Crawler Endpoint history will be erased.
func ResetCollector(spiderID string, bucket store.Bucket) error {

	collector, err := store.NewCollectStorage(spiderID, bucket)
	if err != nil {
		return err
	}

	return collector.Reset()
}

// ResetExtractor drops the extractor store
// Extracted data will be erased. All temporary data/images will be lost.
func ResetExtractor(spiderID string, bucket store.Bucket) error {

	extractor, err := store.NewExtractStorage(spiderID, bucket)
	if err != nil {
		return err
	}

	return extractor.Reset()
}

// ResetMedia drops the media store
// All media files will be erased.
func ResetMedia(spiderID string, bucket store.Bucket) error {

	media, err := store.NewMediaStorage(spiderID, bucket)
	if err != nil {
		return err
	}

	return media.Reset()
}

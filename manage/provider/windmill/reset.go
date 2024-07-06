package windmill

import (
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
//
//goland:noinspection GoDfaNilDereference
func Reset(name string) error {

	var deploy *setup.Deploy

	if err := DeployResource(deploy); err != nil {
		return err
	}

	// todo reset media

	if err := ResetCollector(name, deploy.MongoDSN); err != nil {
		return err
	}

	return ResetExtractor(name, deploy.MongoDSN)
}

// ResetCollector drops the collector store
// Crawler Endpoint history will be erased.
func ResetCollector(name string, dsn string) error {

	collector, err := store.NewCollectStore(name, dsn)
	if err != nil {
		return err
	}

	return collector.Drop()
}

// ResetExtractor drops the extractor store
// Extracted data will be erased. All temporary data/images will be lost.
func ResetExtractor(name string, dsn string) error {

	extractor, err := store.NewExtractStore(name, dsn)
	if err != nil {
		return err
	}

	return extractor.Drop()
}

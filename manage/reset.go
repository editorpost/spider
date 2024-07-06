package manage

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
func Reset(name string, cfg *mongodb.Config) error {

	collector, err := store.NewCollectStore(name, cfg.DSN)
	if err != nil {
		return err
	}

	err = collector.Drop()
	if err != nil {
		return err
	}

	extractor, err := store.NewExtractStore(name, cfg.DSN)
	if err != nil {
		return err
	}

	return extractor.Drop()
}

// ResetCollector drops the collector store
// crawler Endpoint history will be erased.
func ResetCollector(name string, cfg *mongodb.Config) error {

	collector, err := store.NewCollectStore(name, cfg.DSN)
	if err != nil {
		return err
	}

	return collector.Drop()
}

// ResetExtractor drops the extractor store
// Extracted data will be erased. All temporary data/images will be lost.
func ResetExtractor(name string, cfg *mongodb.Config) error {

	extractor, err := store.NewExtractStore(name, cfg.DSN)
	if err != nil {
		return err
	}

	return extractor.Drop()
}

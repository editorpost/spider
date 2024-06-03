package windmill

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
func Reset(name string) error {

	conf, err := MongoConfig(DefaultMongoResource)
	if err != nil {
		return err
	}

	if err := ResetCollector(name, conf); err != nil {
		return err
	}

	return ResetExtractor(name, conf)
}

// ResetCollector drops the collector store
// Crawler URL history will be erased.
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

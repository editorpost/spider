package manage

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/store"
)

// Reset drops the collector and extractor stores
// All spider related data will be erased.
//
//goland:noinspection GoUnusedExportedFunction
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

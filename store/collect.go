package store

import (
	"errors"
	"github.com/gocolly/colly/v2/storage"
)

const (
	CollectorVisited  = "visited"
	CollectorCookies  = "cookies"
	ExtractorResults  = "crawled"
	StoreRequestIDKey = "requestID"
	StoreVisitedKey   = "visited"
	StoreHostKey      = "host"
	StoreCookiesKey   = "cookies"
)

type CollectStore interface {
	storage.Storage
	Drop() error
}

func NewCollectStore(dbName, dsn string) (CollectStore, error) {

	if dsn == "" {
		return nil, errors.New("collector store dsn is empty")
	}

	// switch store by dsn prefix database, e.g. mongodb://, postgres://
	switch {
	case dsn[:8] == "mongodb://":
		return NewMongoCollectStore(dbName, dsn)
	}

	// return error if dsn not supported
	return nil, errors.New("can't create store for dsn: " + dsn + ", only mongodb:// and postgres:// are supported")
}

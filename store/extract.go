package store

import (
	"context"
	"errors"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/extract"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"runtime/debug"
)

type (

	// ExtractStore with deduplication based on unique field
	ExtractStore struct {
		client          *mongo.Client
		db              *mongo.Database
		col             *mongo.Collection
		cache           *bloom.BloomFilter
		uniqueField     string
		uniqueOverwrite bool
	}
)

func NewExtractStore(jobDbName string, cfg *mongodb.Config) (s *ExtractStore, err error) {

	s = &ExtractStore{}
	uri := options.Client().ApplyURI(cfg.DSN)

	if s.client, err = mongo.Connect(context.Background(), uri); err != nil {
		return nil, err
	}

	s.db = s.client.Database(jobDbName)
	s.col = s.db.Collection(ExtractorResults)
	s.cache = bloom.NewWithEstimates(1000000, 0.01)

	// do not uniqueOverwrite by default
	s.uniqueOverwrite = false
	s.uniqueField = extract.UrlField

	return
}

// Save saves extracted data to mongo
func (s *ExtractStore) Save(p *extract.Payload) error {

	if err := s.save(p.Data); err != nil {
		return err
	}

	slog.Info("saved", slog.String("url", p.URL.String()))
	return nil
}

func (s *ExtractStore) Client() *mongo.Client {
	return s.client
}

func (s *ExtractStore) Close() error {
	return s.client.Disconnect(context.Background())
}

func (s *ExtractStore) save(row map[string]any) error {

	if len(s.uniqueField) > 0 {
		return s.upsert(row)
	}

	return s.insert(row)
}

func (s *ExtractStore) upsert(row map[string]any) error {

	// check if item exists
	field := s.uniqueField
	maybeNew := s.absentPerRun(bson.M{field: row[field]})

	if maybeNew {

		// read from mongo
		_, err := s.read(bson.M{field: row[field]})

		// not found means new item
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s.insert(row)
		}

		// internal db error, fail
		if err != nil {
			slog.Error("read db error: ", err)
			return err
		}

		// row exists, no need to overwrite
		if !s.uniqueOverwrite {
			return nil
		}
	}

	// already exists, update
	err := s.update(bson.M{field: row[field]}, row)
	if err != nil {
		slog.Error("update db error: ", err)
	}

	return err
}

func (s *ExtractStore) read(req bson.M) (map[string]any, error) {

	var row map[string]any
	err := s.col.FindOne(context.Background(), req).Decode(&row)
	return row, err
}

// absentPerRun returns true if item is in cache (was processed during this runtime session)
func (s *ExtractStore) absentPerRun(req bson.M) bool {
	return !s.cache.Test([]byte(req[s.uniqueField].(string)))
}

func (s *ExtractStore) insert(row map[string]any) error {

	_, err := s.col.InsertOne(context.Background(), row)

	if err != nil {
		slog.Error("save item error: " + err.Error())
		debug.PrintStack()
	}

	s.cache.Add([]byte(row[s.uniqueField].(string)))

	return err
}

func (s *ExtractStore) update(req bson.M, row map[string]any) error {

	_, err := s.col.ReplaceOne(context.Background(), req, row)
	s.cache.Add([]byte(row[s.uniqueField].(string)))

	return err
}

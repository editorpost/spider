package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"runtime/debug"
	"spider/config"
)

type (
	ExtractOptions struct {
		// Database mongo database name
		Database string
		// URI mongo connection string
		URI string
		// CollectionName mongo collection name to store extracted data
		CollectionName string
		// TaskID crawlab task id
		TaskID string
		// UniqueField field name (e.g. url, id, etc.) for overwriting/deduplication
		UniqueField string
	}

	// ExtractStore implements a MongoDB storage backend for colly
	ExtractStore struct {
		opt    ExtractOptions
		client *mongo.Client
		db     *mongo.Database
		col    *mongo.Collection
	}
)

// NewExtractStoreFromEnv follows env vars conventions
// defined in
func NewExtractStoreFromEnv() *ExtractStore {

	cfg := config.GetEnv()

	return NewExtractStore(ExtractOptions{
		Database:       cfg.DbName,
		URI:            "mongodb://" + cfg.DbHost + ":" + cfg.DbPort,
		CollectionName: cfg.DbCollection,
		TaskID:         cfg.TaskID,
		UniqueField:    cfg.DbUniqueField,
	})
}

func NewExtractStore(opt ExtractOptions) *ExtractStore {

	if opt.URI == "" {
		slog.Info("no mongo uri provided, using default")
		opt.URI = "mongodb://localhost:27018"
	}

	if opt.Database == "" {
		panic("no mongo database provided for extract store")
	}

	if opt.CollectionName == "" {
		panic("no mongo collection provided for extract store")
	}

	store := &ExtractStore{
		opt: opt,
	}

	if err := store.connect(); err != nil {
		panic("can't connect to mongo db")
	}

	return store
}

// connect initializes the MongoDB storage
func (s *ExtractStore) connect() error {

	var err error

	uri := options.
		Client().
		SetHosts([]string{s.opt.URI}).
		ApplyURI(s.opt.URI)

	s.client, err = mongo.Connect(context.Background(), uri)
	if err != nil {
		return err
	}

	s.db = s.client.Database(s.opt.Database)
	s.col = s.db.Collection(s.opt.CollectionName)

	return nil
}

func (s *ExtractStore) Client() *mongo.Client {
	return s.client
}

func (s *ExtractStore) Close() error {
	return s.client.Disconnect(context.Background())
}

func (s *ExtractStore) Save(row map[string]any) error {

	row["_tid"] = s.opt.TaskID

	if len(s.opt.UniqueField) > 0 {
		return s.overwrite(row)
	}

	return s.insert(row)
}

func (s *ExtractStore) overwrite(row map[string]any) error {

	// check if item exists
	field := s.opt.UniqueField
	_, err := s.read(bson.M{field: row[field]})

	// not found
	if errors.Is(err, mongo.ErrNoDocuments) {
		return s.insert(row)
	}

	// db error
	if err != nil {
		slog.Error("read db error: ", err)
		return err
	}

	// already exists
	err = s.update(bson.M{field: row[field]}, row)
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

func (s *ExtractStore) insert(row map[string]any) error {

	_, err := s.col.InsertOne(context.Background(), row)

	if err != nil {
		slog.Error("save item error: " + err.Error())
		debug.PrintStack()
	}

	return err
}

func (s *ExtractStore) update(req bson.M, row map[string]any) error {

	_, err := s.col.UpdateOne(context.Background(), req, row)
	return err
}

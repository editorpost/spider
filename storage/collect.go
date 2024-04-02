package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"net/url"
	"strconv"
)

const (
	StoreVisitedCollection = "colly_visited"
	StoreCookiesCollection = "colly_cookies"
	StoreRequestIDKey      = "requestID"
	StoreVisitedKey        = "visited"
	StoreHostKey           = "host"
	StoreCookiesKey        = "cookies"
)

// CollectStore implements a MongoDB storage backend for colly
type CollectStore struct {
	Database string
	URI      string
	client   *mongo.Client
	db       *mongo.Database
	visited  *mongo.Collection
	cookies  *mongo.Collection
}

func NewCollectorStore(database, uri string) *CollectStore {

	if uri == "" {
		uri = "mongodb://localhost:27018"
	}

	return &CollectStore{
		Database: database,
		URI:      uri,
	}
}

// Init initializes the MongoDB storage
func (s *CollectStore) Init() error {

	var err error

	uri := options.
		Client().
		ApplyURI(s.URI).
		SetAuth(options.Credential{
			Username: "root",
			Password: "nopass",
		})

	s.client, err = mongo.Connect(context.Background(), uri)
	if err != nil {
		return err
	}

	s.db = s.client.Database(s.Database)
	s.visited = s.db.Collection(StoreVisitedCollection)
	s.cookies = s.db.Collection(StoreCookiesCollection)

	return nil
}

// Visited implements colly/storage.Visited()
func (s *CollectStore) Visited(requestID uint64) error {

	_, err := s.visited.InsertOne(context.Background(), bson.D{
		{StoreRequestIDKey, strconv.FormatUint(requestID, 10)},
		{StoreVisitedKey, true},
	})

	return err
}

// IsVisited implements colly/storage.IsVisited()
func (s *CollectStore) IsVisited(requestID uint64) (bool, error) {

	doc := bson.D{}
	err := s.visited.FindOne(nil, bson.D{
		{StoreRequestIDKey, strconv.FormatUint(requestID, 10)},
	}).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Cookies implements colly/storage.Cookies()
func (s *CollectStore) Cookies(u *url.URL) string {

	doc := bson.D{}
	err := s.cookies.FindOne(nil, bson.D{
		{StoreHostKey, u.Host},
	}).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			slog.Warn("cookies not found", slog.String("host", u.Host))
		}
		return ""
	}

	return doc.Map()[StoreCookiesKey].(string)
}

// SetCookies implements colly/storage.SetCookies()
func (s *CollectStore) SetCookies(u *url.URL, cookies string) {

	_, err := s.cookies.InsertOne(nil, bson.D{
		{StoreHostKey, u.Host},
		{StoreCookiesKey, cookies},
	})

	if err != nil {
		slog.Warn("set cookies failed", slog.String("host", u.Host), slog.String("cookies", cookies))
	}
}

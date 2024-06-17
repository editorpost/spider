package store

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/bits-and-blooms/bloom/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"net/url"
	"strconv"
	"sync"
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

// CollectStore implements a MongoDB storage backend for colly
type CollectStore struct {
	client   *mongo.Client
	db       *mongo.Database
	visited  *mongo.Collection
	cookies  *mongo.Collection
	_visited *bloom.BloomFilter
	_cookies *sync.Map
}

// Init satisfy colly/storage.Storage interface
func (s *CollectStore) Init() error {
	return s.preload()
}

func NewCollectStore(dbName, mongoURI string) (s *CollectStore, err error) {

	if mongoURI == "" {
		return nil, errors.New("collector store config is nil")
	}

	s = &CollectStore{}
	uri := options.Client().ApplyURI(mongoURI)

	if s.client, err = mongo.Connect(context.Background(), uri); err != nil {
		return
	}

	s.db = s.client.Database(dbName)
	s.visited = s.db.Collection(CollectorVisited)
	s.cookies = s.db.Collection(CollectorCookies)
	s._cookies = &sync.Map{}

	return
}

// preload loads visited urls from db to cache
func (s *CollectStore) preload() error {

	// init visited collection
	s._visited = bloom.NewWithEstimates(1000000, 0.01)

	// load visited urls from db
	cursor, dbErr := s.visited.Find(context.TODO(), bson.D{})

	if errors.Is(dbErr, mongo.ErrNoDocuments) {
		return nil
	}

	if dbErr != nil {
		return dbErr
	}
	//goland:noinspection GoUnhandledErrorResult
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		doc := bson.D{}
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		requestID, err := strconv.ParseUint(doc.Map()[StoreRequestIDKey].(string), 10, 64)
		if err != nil {
			return err
		}

		s.cacheVisited(requestID)
	}

	return nil
}

// Visited implements colly/storage.Visited()
func (s *CollectStore) Visited(requestID uint64) error {

	if requestID == 0 {
		return errors.New("requestID is zero")
	}

	_, err := s.visited.InsertOne(context.Background(), bson.D{
		{StoreRequestIDKey, strconv.FormatUint(requestID, 10)},
		{StoreVisitedKey, true},
	})

	if err != nil {
		slog.Error("visited failed", slog.Uint64("requestID", requestID), slog.String("err", err.Error()))
		return err
	}

	s.cacheVisited(requestID)

	return nil
}

// IsVisited implements colly/storage.IsVisited()
func (s *CollectStore) IsVisited(requestID uint64) (bool, error) {

	// check cache
	if s.hasVisitedCache(requestID) {
		return true, nil
	}

	doc := bson.D{}
	err := s.visited.FindOne(context.TODO(), bson.D{
		{StoreRequestIDKey, strconv.FormatUint(requestID, 10)},
	}).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	s.cacheVisited(requestID)

	return true, nil
}

// Cookies implements colly/storage.Cookies()
func (s *CollectStore) Cookies(u *url.URL) string {

	// check cache
	if v, ok := s._cookies.Load(u.Host); ok {
		if str, k := v.(string); k {
			return str
		}
		return ""
	}

	// check db
	doc := bson.D{}
	err := s.cookies.FindOne(context.TODO(), bson.D{
		{StoreHostKey, u.Host},
	}).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			slog.Debug("cookies have not exist yet", slog.String("host", u.Host))
		}
		return ""
	}

	str := doc.Map()[StoreCookiesKey].(string)
	s._cookies.Store(u.Host, str)
	slog.Info("Set cookies from db", slog.String("host", u.Host), slog.String("cookies", str))

	return str
}

// SetCookies implements colly/storage.SetCookies()
func (s *CollectStore) SetCookies(u *url.URL, cookies string) {

	_, err := s.cookies.InsertOne(context.TODO(), bson.D{
		{StoreHostKey, u.Host},
		{StoreCookiesKey, cookies},
	})

	if err != nil {
		slog.Warn("set cookies failed", slog.String("host", u.Host), slog.String("cookies", cookies))
		return
	}

	s._cookies.Store(u.Host, cookies)
}

// Drop job database collections from storage with all data
func (s *CollectStore) Drop() error {

	if err := s.visited.Drop(context.Background()); err != nil {
		return err
	}

	if err := s.cookies.Drop(context.Background()); err != nil {
		return err
	}

	return nil
}

// cacheVisited sets requestID as visited to cache
func (s *CollectStore) cacheVisited(requestID uint64) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, requestID)
	s._visited.Add(key)
}

// absentPerRun returns true if item is in cache (was processed during this runtime session)
func (s *CollectStore) hasVisitedCache(requestID uint64) bool {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, requestID)
	return s._visited.Test(key)
}

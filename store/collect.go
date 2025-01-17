package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/editorpost/donq/res"
	"github.com/gocolly/colly/v2/storage"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type CollectStore interface {
	storage.Storage
}

// CollectStorage in-memory colly storage backed by S3
// @see CollectHistory.Init and CollectStorage.Shutdown
type CollectStorage struct {
	b        res.S3
	store    Storage
	filepath string
	// Based on colly storage.InMemoryStorage
	// @source github.com/gocolly/colly/v2@v2.1.1-0.20240327170223-5224b972e22b/storage/storage.go
	visitedURLs map[uint64]bool
	lock        *sync.RWMutex
	jar         *cookiejar.Jar
}

func NewCollectStorage(folder string, b res.S3) (*CollectStorage, func() error, error) {

	jar, _ := cookiejar.New(nil)

	store, err := NewStorage(b, folder)
	if err != nil {
		return nil, nil, fmt.Errorf("extract store s3 client: %w", err)
	}

	s := &CollectStorage{
		b:           b,
		store:       store,
		filepath:    VisitedFile,
		visitedURLs: make(map[uint64]bool),
		lock:        &sync.RWMutex{},
		jar:         jar,
	}

	return s, s.shutdown, nil
}

func DropCollectStorage(folder string, b res.S3) error {

	// note no need to call shutdown here
	store, _, err := NewCollectStorage(folder, b)
	if err != nil {
		return err
	}

	return store.Reset()
}

func (s *CollectStorage) shutdown() error {

	if len(s.visitedURLs) == 0 {
		return nil
	}

	b, err := json.Marshal(s.visitedURLs)
	if err != nil {
		return err
	}

	return s.store.Save(b, s.filepath)
}

// Init initializes CollectStorage
func (s *CollectStorage) Init() error {

	b, err := s.store.Load(s.filepath)

	// Check for the "Not Found" error
	var awsErr *types.NoSuchKey
	if errors.As(err, &awsErr) {
		return nil
	}

	// ignore not found error
	if err != nil {
		return fmt.Errorf("error to connect to collect storage, %w", err)
	}

	return json.Unmarshal(b, &s.visitedURLs)
}

func (s *CollectStorage) Reset() error {
	return s.store.Delete(s.filepath)
}

// Visited implements Storage.Visited()
func (s *CollectStorage) Visited(requestID uint64) error {
	s.lock.Lock()
	s.visitedURLs[requestID] = true
	s.lock.Unlock()
	return nil
}

// IsVisited implements Storage.IsVisited()
func (s *CollectStorage) IsVisited(requestID uint64) (bool, error) {
	s.lock.RLock()
	visited := s.visitedURLs[requestID]
	s.lock.RUnlock()
	return visited, nil
}

// Cookies implements Storage.Cookies()
func (s *CollectStorage) Cookies(u *url.URL) string {
	return storage.StringifyCookies(s.jar.Cookies(u))
}

// SetCookies implements Storage.SetCookies()
func (s *CollectStorage) SetCookies(u *url.URL, cookies string) {
	s.jar.SetCookies(u, storage.UnstringifyCookies(cookies))
}

// Close implements Storage.Close()
func (s *CollectStorage) Close() error {
	return nil
}

package store

import (
	"encoding/json"
	"fmt"
	"github.com/editorpost/donq/res"
	"github.com/gocolly/colly/v2/storage"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type CollectStore interface {
	storage.Storage
	Reset() error
}

// CollectStorage in-memory colly storage backed by S3
// @see CollectHistory.Init and CollectStorage.Shutdown
type CollectStorage struct {
	b     res.S3
	store Storage
	// Based on colly storage.InMemoryStorage
	// @source github.com/gocolly/colly/v2@v2.1.1-0.20240327170223-5224b972e22b/storage/storage.go
	visitedURLs map[uint64]bool
	lock        *sync.RWMutex
	jar         *cookiejar.Jar
}

func NewCollectStorage(spiderID string, b res.S3) (*CollectStorage, func() error, error) {

	jar, _ := cookiejar.New(nil)
	folder := fmt.Sprintf(CollectFolder, spiderID)

	store, err := NewStorage(b, folder)
	if err != nil {
		return nil, nil, fmt.Errorf("extract store s3 client: %w", err)
	}

	s := &CollectStorage{
		b:           b,
		store:       store,
		visitedURLs: make(map[uint64]bool),
		lock:        &sync.RWMutex{},
		jar:         jar,
	}

	return s, s.shutdown, nil
}

func (s *CollectStorage) shutdown() error {

	if len(s.visitedURLs) == 0 {
		return nil
	}

	b, err := json.Marshal(s.visitedURLs)
	if err != nil {
		return err
	}

	return s.store.Save(b, VisitedFile)
}

// Init initializes CollectStorage
func (s *CollectStorage) Init() error {

	b, err := s.store.Load(VisitedFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &s.visitedURLs)
}

func (s *CollectStorage) Reset() error {
	return s.store.Delete(VisitedFile)
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

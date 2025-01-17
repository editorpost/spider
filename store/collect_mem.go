package store

import (
	"github.com/gocolly/colly/v2/storage"
	"hash/fnv"
	"log/slog"
)

// ExtractorHistory in-memory colly storage backed by S3
// @see CollectHistory.Init and ExtractorHistory.Shutdown
type ExtractorHistory struct {
	storage.InMemoryStorage
	// Based on colly storage.InMemoryStorage
	// @source github.com/gocolly/colly/v2@v2.1.1-0.20240327170223-5224b972e22b/storage/storage.go
	visitedURLs map[uint64]bool
}

func NewMemoryCollector(visited func() []string) *ExtractorHistory {

	visitedURLs := make(map[uint64]bool)

	for _, url := range visited() {

		hash, err := FNVHash(url)
		if err != nil {
			slog.Error("extractor history hash", "err", err.Error())
			continue
		}

		visitedURLs[hash] = true
	}

	return &ExtractorHistory{
		visitedURLs: visitedURLs,
	}
}

func FNVHash(uri string) (uint64, error) {

	h := fnv.New64a()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return 0, err
	}

	return h.Sum64(), nil
}

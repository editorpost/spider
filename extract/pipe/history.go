package pipe

import (
	"hash/fnv"
	"log/slog"
	"sync"
)

// History is a store for payloads extraction history.
// Implements configuration flag Collect.ExtractOnce avoiding duplication of payloads.
//
// However, the spider might be configured to extract multiple payloads from the same page.
// So, new payloads are counted as extracted only if they are from different pages.
//
// If URL matched with extracted on previous runs extraction will be skipped.
type (
	History struct {
		visitedURLs map[uint64]uint32
		mute        *sync.RWMutex
	}
)

func NewPayloadHistory() *History {

	history := &History{
		visitedURLs: make(map[uint64]uint32),
		mute:        &sync.RWMutex{},
	}

	return history
}

func (s *History) Init(loader func() []string) {
	for _, u := range loader() {
		s.extracted(0, u)
	}
}

func (s *History) Extracted(p *Payload) {
	s.extracted(p.Doc.Request.ID, p.URL.String())
}

// IsExtracted used as starter pipeline processor
func (s *History) IsExtracted(p *Payload) (bool, error) {
	return s.isExtracted(p.Doc.Request.ID, p.URL.String())
}

// extracted marks a payload as extracted
func (s *History) extracted(requestID uint32, url string) {

	// if colly request id and url are the same,
	// then we can allow multiple payloads per page

	hash, err := FNVHash(url)
	if err != nil {
		slog.Error("payloads extraction history hash", "err", err.Error())
	}

	s.mute.Lock()
	defer s.mute.Unlock()
	s.visitedURLs[hash] = requestID
}

// IsExtracted checks if a URL has been extracted.
// Not count payload as extracted until the request id is the same,
// so it doesn't prevent duplicates caused double function call
// Returns true if the URL has been extracted before.
func (s *History) isExtracted(requestID uint32, url string) (bool, error) {

	hash, err := FNVHash(url)
	if err != nil {
		slog.Error("payloads extraction history hash", "err", err.Error())
	}

	s.mute.RLock()
	defer s.mute.RUnlock()

	reqID, ok := s.visitedURLs[hash]
	if !ok {
		return false, nil
	}

	// Since multiple payloads can be extracted from a single page,
	// we should allow duplicated urls, but during same request id.
	// The requestID is sequential number of the request during colly run.
	//
	// Here we compare current request id with the extracted request id.
	// If they are the same, then we count payload not yet extracted.
	return reqID != requestID, nil
}

func FNVHash(uri string) (uint64, error) {

	h := fnv.New64a()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return 0, err
	}

	return h.Sum64(), nil
}

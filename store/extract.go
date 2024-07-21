package store

import (
	"encoding/json"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/editorpost/spider/extract/pipe"
)

const ()

type (
	ExtractStore interface {
		Save(p *pipe.Payload) error
		Reset() error
		Close() error
	}

	ExtractStorage struct {
		b         Bucket
		store     Storage
		extracted *bloom.BloomFilter
	}
)

// NewExtractStorage S3 storage for payload and document html
func NewExtractStorage(spiderID string, b Bucket) (*ExtractStorage, error) {

	store, err := NewStorage(b, fmt.Sprintf(PayloadFolder, spiderID))
	if err != nil {
		return nil, fmt.Errorf("extract store s3 client: %w", err)
	}

	return &ExtractStorage{
		b:         b,
		store:     store,
		extracted: bloom.NewWithEstimates(100000, 0.01),
	}, nil
}

func (s *ExtractStorage) Load(filename string) ([]byte, error) {
	return s.store.Load(filename)
}

func (s *ExtractStorage) Save(p *pipe.Payload) (err error) {

	key := p.URL.String()

	// bloom: only false values trusted
	if !s.extracted.Test([]byte(key)) {

		if err = s.save(p); err != nil {
			s.extracted.Add([]byte(key))
		}
	}

	return nil
}

func (s *ExtractStorage) save(p *pipe.Payload) error {

	b, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}

	// payload
	err = s.store.Save(b, fmt.Sprintf("%s/%s", p.ID, PayloadFile))
	if err != nil {
		return err
	}

	// html
	dom, err := p.Doc.DOM.Html()
	if err != nil {
		return err
	}

	return s.store.Save([]byte(dom), fmt.Sprintf("%s/%s", p.ID, DocumentFile))
}

func (s *ExtractStorage) Reset() error {
	return s.store.Reset()
}

func (s *ExtractStorage) Close() error {
	return nil
}

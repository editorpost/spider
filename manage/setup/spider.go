package setup

import (
	"encoding/json"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/payload"
)

type Spider struct {
	*config.Args
	*extract.Config
	pipe *payload.Pipeline
}

// UnmarshalJSON is the custom unmarshalling for Spider
func (s *Spider) UnmarshalJSON(data []byte) error {
	type Alias Spider
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func NewSpider(data []byte) (*Spider, error) {

	s := &Spider{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	if err := s.withPipeline(); err != nil {
		return nil, err
	}

	return s, nil
}

// Pipeline for a crawler hooks with extractor functions.
// Might have prepopulated extractors or extractors added later.
// Empty by default.
func (s *Spider) Pipeline() *payload.Pipeline {
	return s.pipe
}

// WithPipeline sets up the extractors for the spider.
func (s *Spider) withPipeline() error {

	if s.pipe != nil {
		return nil
	}

	extractors, err := extract.Extractors(s.ExtractFields, s.ExtractEntities...)
	if err != nil {
		return err
	}

	s.pipe = payload.NewPipeline(extractors...)
	return nil
}

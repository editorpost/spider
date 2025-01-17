package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/media"
)

type Config struct {
	// ExtractEntity is predefined named entity to extract
	Entities []string `json:"Entities"`
	// Fields is the list of fields to extract
	Fields []*fields.Field `json:"Fields"`
	// Media is the configuration for media
	Media *media.Config `json:"Media"`
	// ExtractOnce is the flag to extract the entity only once
	// If true, then existing payloads urls loaded from db
	// If false, payloads are extracted from the page and stored in db without unique check
	ExtractOnce bool `json:"ExtractOnce"`
}

func (c *Config) Normalize() error {

	if c.Media == nil {
		c.Media = &media.Config{}
	}

	if c.Fields == nil {
		c.Fields = make([]*fields.Field, 0)
	}

	return nil
}

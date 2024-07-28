package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/media"
)

type Config struct {
	// ExtractEntity is predefined named entity to extract
	Extract []string `json:"Extract"`
	// Fields is the list of fields to extract
	Fields []*fields.Field `json:"Fields"`
	// Media is the configuration for media
	Media *media.Config `json:"Media"`
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

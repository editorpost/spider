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

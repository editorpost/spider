package extract

import "github.com/editorpost/spider/extract/fields"

type Config struct {
	// ExtractEntity is predefined named entity to extract
	ExtractEntities []string `json:"ExtractEntities"`
	// ExtractFields is the list of fields to extract
	ExtractFields []*fields.Field `json:"ExtractFields"`
}

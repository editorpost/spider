package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/payload"
)

// Fields extracts the fields from the HTML
// and sets the fields to the payload
func Fields(root ...*fields.Field) (payload.Extractor, error) {

	extract, err := fields.Extractor(root...)
	if err != nil {
		return nil, err
	}

	return func(p *payload.Payload) error {

		data := map[string]any{}
		extract(data, p.Selection)

		if len(data) == 0 {
			return payload.ErrDataNotFound
		}

		for k, v := range data {
			p.Data[k] = v
		}

		return nil
		// closure
	}, nil // return
}

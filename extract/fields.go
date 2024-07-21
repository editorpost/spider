package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/pipe"
)

// Fields extracts the fields from the HTML
// and sets the fields to the payload
func Fields(root ...*fields.Field) (pipe.Extractor, error) {

	extract, err := fields.Extractor(root...)
	if err != nil {
		return nil, err
	}

	return func(p *pipe.Payload) error {

		data := map[string]any{}
		extract(data, p.Selection)

		if len(data) == 0 {
			return pipe.ErrDataNotFound
		}

		for k, v := range data {
			p.Data[k] = v
		}

		return nil
		// closure
	}, nil // return
}

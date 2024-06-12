package extract

import (
	"github.com/editorpost/spider/extract/fields"
)

// Fields extracts the fields from the HTML
// and sets the fields to the payload
func Fields(root ...*fields.Field) (PipeFn, error) {

	extract, err := fields.Extractor(root...)
	if err != nil {
		return nil, err
	}

	return func(payload *Payload) error {

		data := map[string]any{}
		extract(data, payload.Selection)

		if len(data) == 0 {
			return ErrDataNotFound
		}

		for k, v := range data {
			payload.Data[k] = v
		}

		return nil
		// closure
	}, nil // return
}

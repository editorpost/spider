package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"log/slog"
)

// Fields extracts the fields from the HTML
// and sets the fields to the payload
func Fields(extractors ...*fields.Extractor) func(*Payload) error {

	callbacks := map[string]fields.ExtractFn{}

	for _, d := range extractors {

		extractor, err := fields.Build(d)
		if err != nil {
			slog.Error("failed to build extractor", slog.String("error", err.Error()))
		}

		callbacks[d.FieldName] = extractor
	}

	return func(p *Payload) error {

		// extract fields
		for _, extractor := range extractors {

			values, err := callbacks[extractor.FieldName](p.Selection)
			if err != nil {
				return err
			}

			// single value
			if extractor.Limit == 1 && len(values) > 0 {
				p.Data[extractor.FieldName] = values[0]
				return nil
			}

			// cut off values if limit is set
			if extractor.Limit > 0 && len(values) > extractor.Limit {
				values = values[:extractor.Limit]
			}
		}

		return nil
	}
}

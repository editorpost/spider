package extract

import (
	"github.com/editorpost/spider/extract/fields"
	"log/slog"
)

// Fields extracts the fields from the HTML
// and sets the fields to the payload
func Fields(extractors ...*fields.Extractor) func(*Payload) error {

	extractFn := FieldsExtractFn(extractors...)

	return func(p *Payload) error {

		// extract fields
		for _, extractor := range extractors {

			// todo part inside must be moved to the fields package
			// todo may be in common Build function
			// maybe group and field extractors must not be decoupled

			values, err := extractFn[extractor.FieldName](p.Selection)
			if err != nil {
				return err
			}

			// todo ErrRequiredFieldMissing check required?

			p.Data[extractor.FieldName] = values
		}

		return nil
	}
}

func FieldsExtractFn(extractors ...*fields.Extractor) map[string]fields.ExtractFn {

	extractFn := map[string]fields.ExtractFn{}

	for _, d := range extractors {

		extractor, err := fields.Build(d)
		if err != nil {
			slog.Error("failed to build extractor", slog.String("error", err.Error()))
		}

		extractFn[d.FieldName] = extractor
	}

	return extractFn
}

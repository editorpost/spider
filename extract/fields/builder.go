package fields

import (
	"github.com/PuerkitoBio/goquery"
)

type (
	ExtractFn func(*goquery.Selection) (map[string]any, error)

	Builder interface {
		GetName() string
		IsRequired() bool
		Extractor() (ExtractFn, error)
	}
)

// Build creates a map of group or field names to their corresponding ExtractFn.
func Build[T Builder](bb ...T) (map[string]ExtractFn, error) {

	fns := map[string]ExtractFn{}

	for _, b := range bb {

		fn, err := b.Extractor()
		if err != nil {
			return nil, err
		}

		fns[b.GetName()] = fn
	}

	return fns, nil
}

// Extract is a function that extracts data from a selection.
func Extract[T Builder](name string, builders ...T) (ExtractFn, error) {

	extractors, initErr := Build(builders...)
	if initErr != nil {
		return nil, initErr
	}

	return func(selection *goquery.Selection) (map[string]any, error) {

		// data is a map of field names to their extracted values
		// max entries for group based on Group or Field Cardinality
		data := map[string]any{}

		for _, builder := range builders {

			// apply each extractor to the selection
			// nil | map[string]string | map[string][]string | map[string]any
			values, err := extractors[builder.GetName()](selection)
			if err != nil {
				return nil, err
			}

			for k, v := range values {
				data[k] = v
			}
		}

		// root level fields
		if name == "" {
			return data, nil
		}

		return map[string]any{name: data}, nil
	}, nil
}

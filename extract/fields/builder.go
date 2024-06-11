package fields

import "github.com/PuerkitoBio/goquery"

type (
	ExtractFn func(*goquery.Selection) (any, error)

	Builder interface {
		GetName() string
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

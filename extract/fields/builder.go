package fields

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
)

func Construct(extractor *Extractor) (err error) {

	if err = valid.Struct(extractor); err != nil {
		return err
	}

	if extractor.between, extractor.final, err = RegexCompile(extractor); err != nil {
		return err
	}

	for _, child := range extractor.Children {
		if err = Construct(child); err != nil {
			return err
		}
	}

	return nil
}

func Extract(payload map[string]any, node *goquery.Selection, extractor *Extractor) error {

	var data []any

	if extractor.Children != nil {

		scope := node
		if extractor.Scoped && extractor.Selector != "" {
			scope = node.Find(extractor.Selector)
		}

		deltas := make([]map[string]any, 0)
		scope.Each(func(i int, selection *goquery.Selection) {

			delta := map[string]any{}

			for _, child := range extractor.Children {
				if exErr := Extract(delta, selection, child); exErr != nil {
					return
				}
			}

			if len(delta) > 0 {
				deltas = append(deltas, delta)
			}
		})

		data = lo.ToAnySlice(deltas)
	} else {
		data = lo.ToAnySlice(Value(extractor, node))
	}

	var err error
	if payload[extractor.Name], err = Normalize(data, extractor); err != nil {
		return err
	}

	return err
}

func Value(field *Extractor, sel *goquery.Selection) []string {

	entries := EntriesAsString(field, sel)

	// if regex defined, apply it
	if field.final != nil || field.between != nil {
		entries = RegexPipes(entries, field.between, field.final)
	}

	entries = EntriesTransform(field, entries)
	entries = EntriesClean(entries)

	return entries
}

func Normalize(entries []any, field *Extractor) (any, error) {

	entries = lo.Filter(entries, func(entry any, i int) bool {
		return entry != nil
	})

	if field.Required && len(entries) == 0 {
		return nil, fmt.Errorf("field %s: %w", field.Name, ErrRequiredFieldMissing)
	}

	return Cardinality(field.Cardinality, lo.ToAnySlice(entries)), nil
}

// Cardinality applies cardinality limits to the input entries.
// It used as a final step in the extraction process to convert entries to actual value or field or group.
func Cardinality(cardinality int, entries []any) any {

	if len(entries) == 0 {
		return nil
	}

	// cut to limit len or return all
	if cardinality > 0 && len(entries) > cardinality {
		entries = entries[:cardinality]
	}

	// if limit is 1 return single value
	if cardinality == 1 {
		return entries[0]
	}

	return entries
}

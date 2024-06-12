package fields

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
)

type (
	ExtractFn func(*goquery.Selection) (any, error)

	Builder interface {
		GetName() string
		IsRequired() bool
		Extractor() (ExtractFn, error)
	}
)

// Build creates a map of group or fExtractor names to their corresponding ExtractFn.
func Build(bb ...*Extractor) (map[string]ExtractFn, error) {

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

// ExtractDEpricated is a function that extracts data from a selection.
func ExtractDEpricated(name string, builders ...*Extractor) (ExtractFn, error) {

	_, initErr := Build(builders...)
	if initErr != nil {
		return nil, initErr
	}

	return func(selection *goquery.Selection) (any, error) {

		// data is a map of fExtractor names to their extracted values
		// max entries for group based on Group or Extractor Cardinality
		data := map[string]any{}

		return map[string]any{name: data}, nil
	}, nil
}

func Construct(extractor *Extractor) (err error) {

	if err = valid.Struct(extractor); err != nil {
		return err
	}

	// regex
	if extractor.Between, extractor.Final, err = RegexCompile(extractor); err != nil {
		return err
	}

	for _, child := range extractor.Children {
		if err = construct(child); err != nil {
			return err
		}
	}

	return nil
}

func Extract(payload map[string]any, node *goquery.Selection, extractor *Extractor) error {

	if extractor.Children == nil {

		data := lo.ToAnySlice(extractor.Fieldx(node))
		// todo: remove FieldValue from extractor.Field and use FieldValue call below
		var err error
		payload[extractor.Name], err = FieldValue(data, extractor)
		return err
	}

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

	var err error
	payload[extractor.Name], err = FieldValue(lo.ToAnySlice(deltas), extractor)

	return err
}

package fields

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
)

func Construct(extractor *Extractor) (err error) {

	if err = valid.Struct(extractor); err != nil {
		return err
	}

	if extractor.Between, extractor.Final, err = RegexCompile(extractor); err != nil {
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
		data = lo.ToAnySlice(extractor.Value(node))
	}

	var err error
	if payload[extractor.Name], err = FieldValue(data, extractor); err != nil {
		return err
	}

	return err
}

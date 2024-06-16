package fields

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
)

func Extractor(fields ...*Field) (func(payload map[string]any, node *goquery.Selection), error) {

	if err := Construct(fields...); err != nil {
		return nil, err
	}

	return func(payload map[string]any, node *goquery.Selection) {
		for _, field := range fields {
			Extract(payload, node, field)
		}
	}, nil
}

func Construct(fields ...*Field) (err error) {

	for _, field := range fields {

		if err = valid.Struct(field); err != nil {
			return err
		}

		if field.between, field.final, err = RegexCompile(field); err != nil {
			return err
		}

		for _, child := range field.Children {
			if err = Construct(child); err != nil {
				return err
			}
		}
	}

	return nil
}

func Extract(payload map[string]any, node *goquery.Selection, field *Field) {

	var data []any

	if field.Children != nil {

		scope := node
		if field.Scoped && field.Selector != "" {
			scope = node.Find(field.Selector)
		}

		deltas := make([]map[string]any, 0)
		scope.Each(func(i int, selection *goquery.Selection) {

			delta := map[string]any{}
			for _, child := range field.Children {

				Extract(delta, selection, child)
				if delta[child.Name] == nil && child.Required {
					return
				}
			}

			deltas = append(deltas, delta)
		})

		data = lo.ToAnySlice(deltas)
	} else {
		data = lo.ToAnySlice(Value(field, node))
	}

	values := Normalize(data, field.Cardinality)

	if values != nil {
		payload[field.Name] = values
	}
}

func Value(field *Field, sel *goquery.Selection) []string {

	entries := SelectionsAsStrings(field, sel)

	if field.final != nil || field.between != nil {
		entries = RegexExtracts(entries, field.between, field.final)
	}

	entries = FormatValues(field, entries)
	entries = CleanStrings(entries)

	return entries
}

func Normalize(entries []any, cardinality int) any {

	entries = lo.Filter(entries, func(entry any, i int) bool {
		return entry != nil
	})

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

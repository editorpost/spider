package fields

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/samber/lo"
)

// Group provides data describing custom data extraction for grouped
type Group struct {
	// Name is a key to store the extracted data.
	// required
	Name string `json:"Name" validate:"required"`

	Limit int `json:"Cardinality"`

	// Selector is a CSS selector to find the element for the group.
	// required
	Selector string `json:"Selector" validate:"required"`

	// Required is a flag to check if at least one value is required.
	Required bool `json:"Required"`

	// Fields is a map of sub-field names to their corresponding Field configurations.
	// required
	Fields []*Field `json:"Fields" validate:"required,dive,required"`

	extractors map[string]ExtractFn
}

// Extractor in case of group, fields extracted by selection
// every extractor has own limited selection area (OuterHtml).
// Result is a slice of maps with extracted
func (group *Group) Extractor() (ExtractFn, error) {

	var e error

	if e = valid.Struct(group); e != nil {
		return nil, e
	}

	if group.extractors, e = Build(group.Fields...); e != nil {
		return nil, e
	}

	extract, e := Extract(group.Fields...)
	if e != nil {
		return nil, e
	}

	// selection might be entity selection or whole document
	return func(selection *goquery.Selection) (any, error) {

		// entries/deltas of the group
		var entries []any

		// select each group entry
		selection.Find(group.Selector).Each(func(i int, groupSelection *goquery.Selection) {
			// entry is a map of field names to their extracted values
			if entry, err := extract(groupSelection); err == nil {
				entries = append(entries, entry)
			}
		})

		if group.Required && len(entries) == 0 {
			return nil, ErrRequiredFieldMissing
		}

		return ApplyCardinality(group.Limit, lo.ToAnySlice(entries)), nil
	}, nil
}

func (group *Group) Map() map[string]any {
	return map[string]any{
		"Name":     group.Name,
		"Selector": group.Selector,
		"Required": group.Required,
		"Fields":   group.Fields,
	}
}

func (group *Group) GetName() string {
	return group.Name
}

func (group *Group) IsRequired() bool {
	return group.Required
}

func GroupFromMap(m map[string]any) (*Group, error) {

	e := &Group{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

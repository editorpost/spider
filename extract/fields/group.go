package fields

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
)

// Group provides data describing custom data extraction for grouped
type Group struct {
	// Name is a key to store the extracted data.
	// required
	Name string `json:"Name" validate:"required"`

	Limit int `json:"Limit"`

	// Selector is a CSS selector to find the element for the group.
	// required
	Selector string `json:"Selector" validate:"required"`

	// Required is a flag to check if at least one value is required.
	Required bool `json:"Required"`

	// Fields is a map of sub-field names to their corresponding Field configurations.
	// required
	Fields map[string]*Field `json:"Fields" validate:"required,dive,required"`
}

// Extractor in case of group, fields extracted by selection
// every extractor has own limited selection area (OuterHtml).
// Result is a slice of maps with extracted
func (group *Group) Extractor() (ExtractFn, error) {

	if err := valid.Struct(group); err != nil {
		return nil, err
	}

	return func(sel *goquery.Selection) (any, error) {

		var values []map[string]any

		// in case of group, fields extracted by selection
		// every extractor has own limited selection area (OuterHtml)

		sel.Find(group.Selector).Each(func(i int, s *goquery.Selection) {
			groupData := make(map[string]any)

			for fieldName, extractor := range group.Fields {

				extractFn, err := extractor.Extractor()

				// stop group extraction
				// error catch under extract.Pipe handler
				if errors.Is(err, ErrRequiredFieldMissing) {
					return
				}

				if err != nil {
					continue
				}

				// clean, unique entries
				entries, err := extractFn(s)
				if err != nil {
					continue
				}

				groupData[fieldName] = entries

			}
			values = append(values, groupData)
		})

		if group.Required && len(values) == 0 {
			return nil, ErrRequiredFieldMissing
		}

		if group.Limit > 0 && len(values) > group.Limit {
			values = values[:group.Limit]
		}

		if group.Limit == 1 {
			return values[0], nil
		}

		return values, nil
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

func GroupFromMap(m map[string]any) (*Group, error) {

	e := &Group{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

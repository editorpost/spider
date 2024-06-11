package fields

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/samber/lo"
	"log/slog"
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
	Fields []*Field `json:"Fields" validate:"required,dive,required"`

	extract map[string]ExtractFn
}

// Extractor in case of group, fields extracted by selection
// every extractor has own limited selection area (OuterHtml).
// Result is a slice of maps with extracted
func (group *Group) Extractor() (ExtractFn, error) {

	var e error

	if e = valid.Struct(group); e != nil {
		return nil, e
	}

	if group.extract, e = ExtractorMap(group.Fields...); e != nil {
		return nil, e
	}

	// selection might be entity selection or whole document
	return func(selection *goquery.Selection) (any, error) {

		// entries/deltas of the group
		var entries []map[string]any

		// select each group entry
		selection.Find(group.Selector).Each(func(i int, groupSelection *goquery.Selection) {
			// entry is a map of field names to their extracted values
			if entry, err := group.EntryFromSelection(groupSelection); err == nil {
				entries = append(entries, entry)
			}
		})

		if group.Required && len(entries) == 0 {
			return nil, ErrRequiredFieldMissing
		}

		return ApplyCardinality(group.Limit, lo.ToAnySlice(entries)), nil
	}, nil
}

func (group *Group) EntryFromSelection(selection *goquery.Selection) (map[string]any, error) {

	// entry is a map of field names to their extracted values
	// max entries for group based on Group.Cardinality
	entry := make(map[string]any)

	// in group selection extract each field
	for _, field := range group.Fields {

		// nil, string, []string, map[string]any
		fieldValue, err := group.extract[field.FieldName](selection)

		// skip group selection if required field is missing
		if errors.Is(err, ErrRequiredFieldMissing) {
			return nil, err
		}

		// log extraction errors
		if err != nil {
			slog.Warn("group field extraction error", "field", field.FieldName, "error", err.Error())
			continue
		}

		if fieldValue == nil {
			continue
		}

		entry[field.FieldName] = fieldValue
	}

	return entry, nil
}

func (group *Group) NormalizeEntries(entries []map[string]any) any {

	if group.Limit > 0 && len(entries) > group.Limit {
		entries = entries[:group.Limit]
	}

	if group.Limit == 1 {
		return entries[0]
	}

	return entries
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

func ExtractorMap(fields ...*Field) (map[string]ExtractFn, error) {

	fns := map[string]ExtractFn{}

	for _, field := range fields {

		fn, err := field.Extractor()
		if err != nil {
			return nil, err
		}

		fns[field.FieldName] = fn
	}

	return fns, nil
}

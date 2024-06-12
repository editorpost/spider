package fields

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/samber/lo"
	"regexp"
)

var (
	// ErrRequiredFieldMissing expected error, stops the pipe chain
	// used by extractors to quickly skip the pipe chain
	// if data is not satisfied or exists
	ErrRequiredFieldMissing = errors.New("skip entity extraction, required field is missing")
)

// Extractor provides data describing custom data extraction from text or html.
type Extractor struct {

	// Name is a key to store the extracted data.
	// required
	Name string `json:"Name" validate:"required"`

	// Cardinality is a number of elements to extract.
	// Zero means no limit, data stored as a list.
	// def: 1
	Cardinality int `json:"Cardinality"`

	// Required is a flag to check if the field is required.
	// If value is falsy then pipeline skips entire entity extraction.
	Required bool `json:"Required"`

	// All formatters apply clear string methods to in/out data:
	// Double spaces are deleted. Output left/right spaces are trimmed.

	// InputFormat is a format of the input data to extractor.
	// It can be "text" or "html".
	// def: "html"
	InputFormat string `json:"InputFormat"`

	// OutputFormat is a format of the output data from extractor.
	// It can be a slice of types "text", "html", "json".
	// Formatters called in the order of the list.
	// def: ["text"]
	OutputFormat []string `json:"OutputFormat"`

	// Selector is a css selector to find the element or limit area for Between/regex.
	// optional
	Selector string `json:"Selector"`

	// Between is a pair of strings to find the element.
	// In case if Selector is not provided, Between applied on whole codfield.

	// BetweenStart is a string to find the element.
	// optional, required if BetweenEnd is provided
	BetweenStart string `json:"BetweenStart"`
	// BetweenEnd is a string to find the element.
	// optional, required if BetweenStart is provided
	BetweenEnd string `json:"BetweenEnd"`

	// FinalRegex is a regular expression to find the element.
	// In case if Selector is not provided, FinalRegex applied on whole codfield.
	// optional
	FinalRegex string `json:"FinalRegex"`

	// Multiline flag prevent deleting new lines from result.
	// optional
	Multiline bool `json:"Multiline"`

	// Scoped flag limits the selection area for the group of field value extractors.
	Scoped bool `json:"Scoped"`

	Between *regexp.Regexp
	Final   *regexp.Regexp

	// Children is a map of sub-field names to their corresponding Extractor configurations.
	// required
	Children []*Extractor `json:"Children" validate:"optional,dive"`

	extract ExtractFn
}

func (field *Extractor) Value(sel *goquery.Selection) []string {

	entries := EntriesAsString(field, sel)

	// if regex defined, apply it
	if field.Final != nil || field.Between != nil {
		entries = RegexPipes(entries, field.Between, field.Final)
	}

	entries = EntriesTransform(field, entries)
	entries = EntriesClean(entries)

	return entries
}

func FieldFromMap(m map[string]any) (*Extractor, error) {

	e := &Extractor{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (field *Extractor) Map() map[string]any {
	return map[string]any{
		"Name":         field.Name,
		"Cardinality":  field.Cardinality,
		"Required":     field.Required,
		"InputFormat":  field.InputFormat,
		"OutputFormat": field.OutputFormat,
		"Selector":     field.Selector,
		"BetweenStart": field.BetweenStart,
		"BetweenEnd":   field.BetweenEnd,
		"FinalRegex":   field.FinalRegex,
		"Multiline":    field.Multiline,
		"Children":     field.Children,
	}
}

func FieldValue(entries []any, field *Extractor) (any, error) {

	entries = lo.Filter(entries, func(entry any, i int) bool {
		return entry != nil
	})

	if field.Required && len(entries) == 0 {
		return nil, fmt.Errorf("field %s: %w", field.Name, ErrRequiredFieldMissing)
	}

	return ApplyCardinality(field.Cardinality, lo.ToAnySlice(entries)), nil
}

// ApplyCardinality applies cardinality limits to the input entries.
// It used as a Final step in the extraction process to convert entries to actual value or field or group.
func ApplyCardinality(cardinality int, entries []any) any {

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

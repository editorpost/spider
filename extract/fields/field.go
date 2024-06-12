package fields

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/samber/lo"
	"regexp"
)

var (
	// ErrRequiredFieldMissing expected error, stops the pipe chain
	// used by extractors to quickly skip the pipe chain
	// if data is not satisfied or exists
	// ErrRequiredFieldMissing
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

// Extractor constructs an extraction function based on the Extractor configuration.
// It validates the Extractor, compiles any necessary regular expressions, and returns a function
// that performs the extraction on a goquery.Selection.
//
// Parameters:
//   - f (*Extractor): A pointer to an Extractor struct containing the configuration for extraction.
//
// Returns:
//   - ExtractFn: A function that takes a goquery.Selection and returns a slice of extracted values or an error.
//   - error: An error that occurred during validation or regex compilation, or nil if successful.
//
// Example:
//
//	extractor := &Extractor{
//	    Name:   "example",
//	    InputFormat: "html",
//	    Selector:    "p",
//	    Cardinality:       2,
//	}
//	extractFn, err := Builder(extractor)
//	if err != nil {
//	    log.Fatalf("Failed to build extractor: %v", err)
//	}
//	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<div><p>Hello</p><p>world!</p></div>"))
//	results, err := extractFn(doc.Selection)
//	if err != nil {
//	    log.Fatalf("Extraction error: %v", err)
//	}
//	fmt.Println(results) // Output: ["Hello", "world!"]
func (field *Extractor) Extractor() (ExtractFn, error) {

	if err := valid.Struct(field); err != nil {
		return nil, err
	}

	var reErr error

	if field.Children != nil {

		// groups requirement
		if field.extract, reErr = ExtractDEpricated("", field.Children...); reErr != nil {
			return nil, reErr
		}

		return field.extractGroup, nil
	}

	// single requirement
	if field.Between, field.Final, reErr = RegexCompile(field); reErr != nil {
		return nil, reErr
	}

	return field.Field, nil
}

func construct(extractor *Extractor) (err error) {

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

func extract(payload map[string]any, node *goquery.Selection, extractor *Extractor) (err error) {

	if extractor.Children == nil {
		payload[extractor.Name], err = extractor.Field(node)
		return err
	}

	scope := node
	if extractor.Scoped && extractor.Selector != "" {
		scope = node.Find(extractor.Selector)
	}

	for _, child := range extractor.Children {
		if err = extract(payload, scope, child); err != nil {
			return err
		}
	}

	return nil
}

// Extractor in case of group, fields extracted by selection
// every extractor has own limited selection area (OuterHtml).
// Result is a slice of maps with extracted
func (field *Extractor) extractGroup(selection *goquery.Selection) (any, error) {

	// entries/deltas of the group
	var entries []any

	selected := selection

	if field.Scoped && field.Selector != "" {
		selected = selection.Find(field.Selector)
	}

	selected.Each(func(i int, groupSelection *goquery.Selection) {
		// entry is a map of fExtractor names to their extracted values
		if entry, err := field.extract(groupSelection); err == nil {
			entries = append(entries, entry)
		}
	})

	return FieldValue(entries, field)
}

func (field *Extractor) Fieldx(sel *goquery.Selection) []string {

	entries := EntriesAsString(field, sel)

	// if regex defined, apply it
	if field.Final != nil || field.Between != nil {
		entries = RegexPipes(entries, field.Between, field.Final)
	}

	entries = EntriesTransform(field, entries)
	entries = EntriesClean(entries)

	return entries
	//return FieldValue(lo.ToAnySlice(entries), field)
}
func (field *Extractor) Field(sel *goquery.Selection) (any, error) {

	entries := EntriesAsString(field, sel)

	// if regex defined, apply it
	if field.Final != nil || field.Between != nil {
		entries = RegexPipes(entries, field.Between, field.Final)
	}

	entries = EntriesTransform(field, entries)
	entries = EntriesClean(entries)

	return FieldValue(lo.ToAnySlice(entries), field)
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

func (field *Extractor) GetName() string {
	return field.Name
}

func (field *Extractor) IsRequired() bool {
	return field.Required
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

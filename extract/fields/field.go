package fields

import (
	"errors"
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

// Field provides data describing custom data extraction from text or html.
type Field struct {

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

	// Selector is a css selector to find the element or limit area for between/regex.
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

	between *regexp.Regexp
	final   *regexp.Regexp
}

// Extractor constructs an extraction function based on the Field configuration.
// It validates the Field, compiles any necessary regular expressions, and returns a function
// that performs the extraction on a goquery.Selection.
//
// Parameters:
//   - f (*Field): A pointer to an Field struct containing the configuration for extraction.
//
// Returns:
//   - ExtractFn: A function that takes a goquery.Selection and returns a slice of extracted values or an error.
//   - error: An error that occurred during validation or regex compilation, or nil if successful.
//
// Example:
//
//	extractor := &Field{
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
func (field *Field) Extractor() (ExtractFn, error) {

	if err := valid.Struct(field); err != nil {
		return nil, err
	}

	var reErr error

	if field.between, field.final, reErr = RegexCompile(field); reErr != nil {
		return nil, reErr
	}

	return field.extract, nil
}

func (field *Field) extract(sel *goquery.Selection) (map[string]any, error) {

	entries := EntriesAsString(field, sel)

	// if regex defined, apply it
	if field.final != nil || field.between != nil {
		entries = RegexPipes(entries, field.between, field.final)
	}

	entries = EntriesTransform(field, entries)
	entries = EntriesClean(entries)

	if field.Required && len(entries) == 0 {
		return nil, ErrRequiredFieldMissing
	}

	return map[string]any{
		field.Name: ApplyCardinality(field.Cardinality, lo.ToAnySlice(entries)),
	}, nil
}

func FieldFromMap(m map[string]any) (*Field, error) {

	e := &Field{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (field *Field) Map() map[string]any {
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
	}
}

func (field *Field) GetName() string {
	return field.Name
}

func (field *Field) IsRequired() bool {
	return field.Required
}

// ApplyCardinality applies cardinality limits to the input entries.
// It used as a final step in the extraction process to convert entries to actual value or field or group.
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

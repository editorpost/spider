package fields

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
)

var (
	// ErrRequiredFieldMissing is an error indicating that the extraction should be skipped
	// because a required field is missing.
	ErrRequiredFieldMissing = errors.New("skip entity extraction, required field is missing")
)

// Field provides data describing custom data extraction from text or html.
type Field struct {

	// FieldName is a key to store the extracted data.
	// required
	FieldName string `json:"FieldName" validate:"required"`

	// Limit is a number of elements to extract.
	// Zero means no limit, data stored as a list.
	// def: 1
	Limit int `json:"Limit"`

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
//	    FieldName:   "example",
//	    InputFormat: "html",
//	    Selector:    "p",
//	    Limit:       2,
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

	between, final, compileErr := RegexCompile(field)
	if compileErr != nil {
		return nil, compileErr
	}

	hasRegex := final != nil || between != nil

	return func(sel *goquery.Selection) (any, error) {

		// css selector selection
		// most time is a single element
		// except cases of parsing listings)
		entries := EntriesAsString(field, sel)

		if hasRegex {
			// most time it is not used
			// but can be useful for complex cases
			//
			// like parsing a list of items
			// it might multiply count of entries
			entries = RegexPipes(entries, between, final)
		}

		// apply output transformers
		entries = EntriesTransform(field, entries)

		// remove duplicates and empty entries
		entries = EntriesClean(entries)

		// if empty field is required
		// skip entire entity extraction
		if len(entries) == 0 {
			if field.Required {
				return nil, ErrRequiredFieldMissing
			}
			return nil, nil
		}

		// final cut to limit len or return all
		if field.Limit > 0 && len(entries) > field.Limit {
			entries = entries[:field.Limit]
		}

		// if limit is 1 return single value
		if field.Limit == 1 {
			return entries[0], nil
		}

		return entries, nil
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
		"FieldName":    field.FieldName,
		"Limit":        field.Limit,
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

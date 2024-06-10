package fields

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
	"strings"
)

var (
	// ErrRequiredFieldMissing is an error indicating that the extraction should be skipped
	// because a required field is missing.
	ErrRequiredFieldMissing = errors.New("skip entity extraction, required field is missing")
)

type (
	ExtractFn func(*goquery.Selection) ([]any, error)

	// UI schema for the extractors:
	// -------------------------------
	// Actions: [Add extractor]
	// -------------------------------
	// New extractor form:
	// -------------------------------
	// FieldName: [article_author]
	// Limit: [1]
	// InputFormat: [html]
	// OutputFormat: [text]
	// Selector: [.author]
	// Between: [start text] [end text]
	// FinalRegex: [by\s+(.+?)\son] // string between `by` and `on`
	// -------------------------------

	// Extractor provides data describing custom data extraction from text or html.
	Extractor struct {

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
		// In case if Selector is not provided, Between applied on whole code.

		// BetweenStart is a string to find the element.
		// optional, required if BetweenEnd is provided
		BetweenStart string `json:"BetweenStart"`
		// BetweenEnd is a string to find the element.
		// optional, required if BetweenStart is provided
		BetweenEnd string `json:"BetweenEnd"`

		// FinalRegex is a regular expression to find the element.
		// In case if Selector is not provided, FinalRegex applied on whole code.
		// optional
		FinalRegex string `json:"FinalRegex"`

		// Multiline flag prevent deleting new lines from result.
		// optional
		Multiline bool `json:"Multiline"`
	}
)

// Build constructs an extraction function based on the Extractor configuration.
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
//	    FieldName:   "example",
//	    InputFormat: "html",
//	    Selector:    "p",
//	    Limit:       2,
//	}
//	extractFn, err := Build(extractor)
//	if err != nil {
//	    log.Fatalf("Failed to build extractor: %v", err)
//	}
//	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<div><p>Hello</p><p>world!</p></div>"))
//	results, err := extractFn(doc.Selection)
//	if err != nil {
//	    log.Fatalf("Extraction error: %v", err)
//	}
//	fmt.Println(results) // Output: ["Hello", "world!"]
func Build(f *Extractor) (ExtractFn, error) {

	if err := valid.Struct(f); err != nil {
		return nil, err
	}

	between, final, compileErr := RegexCompile(f)
	if compileErr != nil {
		return nil, compileErr
	}

	hasRegex := final != nil || between != nil

	return func(sel *goquery.Selection) ([]any, error) {

		// css selector selection
		// most time is a single element
		// except cases of parsing listings)
		entries := EntriesAsString(f, sel)

		if hasRegex {
			// most time it is not used
			// but can be useful for complex cases
			//
			// like parsing a list of items
			// it might multiply count of entries
			entries = RegexPipes(entries, between, final)
		}

		// apply output transformers
		entries = EntriesTransform(f, entries)

		// remove duplicates and empty entries
		entries = EntriesClean(entries)

		// if empty field is required
		// skip entire entity extraction
		if f.Required && len(entries) == 0 {
			return nil, ErrRequiredFieldMissing
		}

		// final cut to limit len or return all
		if f.Limit > 0 && len(entries) > f.Limit {
			entries = entries[:f.Limit]
		}

		return lo.ToAnySlice(entries), nil
	}, nil
}

// EntriesAsString extracts text or HTML content from a goquery.Selection based on the Extractor configuration.
// It processes the selection using the specified input format and optional CSS selector.
//
// Parameters:
//   - f (*Extractor): A pointer to an Extractor struct containing the extraction configuration.
//   - sel (*goquery.Selection): A goquery.Selection object representing the HTML elements to be processed.
//
// Returns:
//   - []string: A slice of strings containing the extracted and processed content.
//
// Example:
//
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<div><p>Hello</p><p>world!</p></div>"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	extractor := &Extractor{
//	    InputFormat: "text",
//	    Selector:    "p",
//	}
//	results := EntriesAsString(extractor, doc.Selection)
//	fmt.Println(results) // Output: ["Hello", "world!"]
func EntriesAsString(f *Extractor, sel *goquery.Selection) []string {

	selection := sel

	if f.Selector != "" {
		// from custom selector
		selection = sel.Find(f.Selector)
	}

	var data []string

	out := func(s *goquery.Selection) string {
		if f.InputFormat == "html" {
			h, _ := goquery.OuterHtml(s)
			return h
		}
		return s.Text()
	}

	selection.Each(func(i int, s *goquery.Selection) {
		data = append(data, ReduceSpaces(out(s)))
	})

	return data
}

// EntriesClean removes empty entries and duplicate entries from the input slice.
// It returns a cleaned slice with only unique, non-empty strings.
//
// Parameters:
//   - entries ([]string): A slice of strings to be cleaned.
//
// Returns:
//   - []string: A cleaned slice of strings containing only unique, non-empty entries,
//     or nil if the cleaned slice is empty.
//
// Example:
//
//	inputs := []string{"apple", "", "banana", "apple", "cherry", ""}
//	cleaned := EntriesClean(inputs)
//	fmt.Println(cleaned) // Output: ["apple", "banana", "cherry"]
func EntriesClean(entries []string) []string {

	// remove empty entries
	entries = lo.Filter(entries, func(v string, i int) bool {
		return v != ""
	})

	entries = lo.Uniq(entries)

	if len(entries) == 0 {
		return nil
	}

	return entries
}

// EntriesTransform applies the EntryTransform function to each entry in the input slice.
// It processes each string according to the Extractor configuration.
//
// Parameters:
//   - f (*Extractor): A pointer to an Extractor struct containing the transformation configuration.
//   - entries ([]string): A slice of strings to be transformed.
//
// Returns:
//   - []string: A slice of transformed strings after applying the EntryTransform function to each entry.
//
// Example:
//
//	extractor := &Extractor{
//	    InputFormat:  "html",
//	    OutputFormat: []string{"text"},
//	}
//	inputs := []string{"<div>Hello  world!</div>", "<p>Go is  awesome!</p>"}
//	outputs := EntriesTransform(extractor, inputs)
//	fmt.Println(outputs) // Output: ["Hello world!", "Go is awesome!"]
func EntriesTransform(f *Extractor, entries []string) []string {

	for i := range entries {
		entries[i] = EntryTransform(f, entries[i])
	}

	return entries
}

// EntryTransform transforms the input value based on the Extractor configuration.
// It processes the value according to the specified input and output formats,
// and applies string manipulation to clean up the output.
//
// Parameters:
//   - f (*Extractor): A pointer to an Extractor struct containing the transformation configuration.
//   - value (string): The input string to be transformed.
//
// Returns:
//   - string: The transformed string after applying the specified input and output format transformations and cleaning.
//
// Example:
//
//	extractor := &Extractor{
//	    InputFormat:  "html",
//	    OutputFormat: []string{"text"},
//	}
//	input := "<div>Hello  world!</div>"
//	output := EntryTransform(extractor, input)
//	fmt.Println(output) // Output: "Hello world!"
func EntryTransform(f *Extractor, value string) string {

	if value == "" {
		return ""
	}

	// apply output transformers aware of the input format
	for _, format := range f.OutputFormat {

		switch format {
		case "text":
			if f.InputFormat == "html" {
				// parse html
				fromHTML, err := goquery.NewDocumentFromReader(strings.NewReader(value))
				if err != nil {
					return ""
				}
				// and extract text
				value = fromHTML.Text()
			}

		case "html":
			// do nothing, value is already html or text
		}
	}

	return ReduceSpaces(value)
}

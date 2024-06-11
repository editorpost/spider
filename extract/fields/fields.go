package fields

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/pkg/valid"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
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

// GroupExtractor provides data describing custom data extraction for grouped
type GroupExtractor struct {
	// Name is a key to store the extracted data.
	// required
	Name string `json:"Name" validate:"required"`

	Limit int `json:"Limit"`

	// Selector is a CSS selector to find the element for the group.
	// required
	Selector string `json:"Selector" validate:"required"`

	// Required is a flag to check if at least one value is required.
	Required bool `json:"Required"`

	// Extractors is a map of sub-field names to their corresponding Extractor configurations.
	// required
	Extractors map[string]*Extractor `json:"Extractors" validate:"required,dive,required"`
}

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
func Build(f any) (ExtractFn, error) {

	switch v := f.(type) {
	case *Extractor:
		return BuildExtractor(v)
	case *GroupExtractor:
		return BuildGroup(v)
	default:
		return nil, fmt.Errorf("unsupported extractor type")
	}
}

func BuildExtractor(f *Extractor) (ExtractFn, error) {

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

func TestBuildExtractor(t *testing.T) {

	tc := []struct {
		name      string
		extractor *Extractor
		expected  []any
		hasErr    bool
		err       error
	}{
		{
			"empty",
			&Extractor{},
			nil,
			true, // field name is required
			nil,
		},
		{
			"simple",
			&Extractor{
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__title",
			},
			[]any{"Main Product Title"},
			false,
			nil,
		},
		{
			"between",
			&Extractor{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__price",
				BetweenStart: "Price:",
				BetweenEnd:   "USD",
			},
			[]any{"99.99"},
			false,
			nil,
		},
		{
			"between image from item prop",
			&Extractor{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"between multiple selections",
			&Extractor{
				FieldName:    "muliple image",
				Limit:        10,
				InputFormat:  "html",
				OutputFormat: []string{"html"},
				Selector:     "meta", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"regex",
			&Extractor{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product--full",
				// multiline regex
				FinalRegex: "Category:(?s)(.*?)</p>",
			},
			[]any{"Magic wands"},
			false,
			nil,
		},
		{
			"regex image from item prop",
			&Extractor{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head",
				FinalRegex:   "meta itemprop=\"image\" content=\"(.+?)\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"all prices",
			&Extractor{
				FieldName: "prices",
				Limit:     0,
				Selector:  ".product__price--amount",
			},
			[]any{"99.99", "49.99", "0.99"},
			false,
			nil,
		},
		{
			"all prices with limit",
			&Extractor{
				FieldName: "prices",
				Limit:     2,
				Selector:  ".product__price--amount",
			},
			[]any{"99.99", "49.99"},
			false,
			nil,
		},
		{
			"required field are empty",
			&Extractor{
				FieldName: "not-exists",
				Limit:     0,
				Selector:  ".product__not-exists-element",
				Required:  true,
			},
			nil,
			true,
			ErrRequiredFieldMissing,
		},
	}

	// use testify assert
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {

			// check error
			skipExpectedErr := func(actual error) bool {

				if actual == nil {
					// continue test case execution
					return false
				}

				// force error if not expected
				if !c.hasErr {
					assert.NoError(t, actual)
				}

				// check error instance
				if c.err != nil {
					assert.ErrorIs(t, c.err, actual)
				}

				// stops test case execution
				return true
			}

			fn, err := Build(c.extractor)
			if skip := skipExpectedErr(err); skip {
				return
			}

			// compare values
			read := strings.NewReader(GetTestFieldsHTML(t))
			dom, err := goquery.NewDocumentFromReader(read)
			require.NoError(t, err)

			values, err := fn(dom.Selection)
			if skip := skipExpectedErr(err); skip {
				return
			}

			assert.Equal(t, c.expected, values)
		})
	}

}

// BuildGroup in case of group, fields extracted by selection
// every extractor has own limited selection area (OuterHtml).
// Result is a slice of maps with extracted
func BuildGroup(f *GroupExtractor) (ExtractFn, error) {

	if err := valid.Struct(f); err != nil {
		return nil, err
	}

	return func(sel *goquery.Selection) ([]any, error) {

		values := []map[string]any{}

		// in case of group, fields extracted by selection
		// every extractor has own limited selection area (OuterHtml)

		sel.Find(f.Selector).Each(func(i int, s *goquery.Selection) {
			groupData := make(map[string]any)

			for fieldName, extractor := range f.Extractors {

				extractFn, err := BuildExtractor(extractor)

				// stop group extraction
				// todo: note this works only under extract.Pipe handler
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

		if f.Required && len(values) == 0 {
			return nil, ErrRequiredFieldMissing
		}

		if f.Limit > 0 && len(values) > f.Limit {
			values = values[:f.Limit]
		}

		return lo.ToAnySlice(values), nil
	}, nil
}

func ExtractorFromMap(m map[string]any) (*Extractor, error) {

	e := &Extractor{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

func TestExtractorFromMap(t *testing.T) {

	m := map[string]any{
		"FieldName":    "title",
		"Limit":        1,
		"InputFormat":  "html",
		"OutputFormat": []string{"text"},
		"Selector":     ".product__title",
	}

	e, err := ExtractorFromMap(m)
	require.NoError(t, err)

	assert.Equal(t, "title", e.FieldName)
	assert.Equal(t, 1, e.Limit)
	assert.Equal(t, "html", e.InputFormat)
	assert.Equal(t, []string{"text"}, e.OutputFormat)
	assert.Equal(t, ".product__title", e.Selector)
}

func GroupFromMap(m map[string]any) (*GroupExtractor, error) {

	e := &GroupExtractor{}
	if err := vars.FromJSON(m, e); err != nil {
		return nil, err
	}

	return e, nil
}

func TestGroupFromMap(t *testing.T) {

	m := map[string]any{
		"Name":     "product",
		"Selector": ".product--full",
		"Required": true,
		"Extractors": map[string]*Extractor{
			"title": {
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__title",
			},
			"price": {
				FieldName:    "price",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__price--amount",
			},
		},
	}

	e, err := GroupFromMap(m)
	require.NoError(t, err)

	assert.Equal(t, "product", e.Name)
	assert.Equal(t, ".product--full", e.Selector)
	assert.True(t, e.Required)
	assert.Len(t, e.Extractors, 2)

	title := e.Extractors["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := e.Extractors["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func TestBuildGroup(t *testing.T) {

	tc := []struct {
		name     string
		group    *GroupExtractor
		expected []map[string]any
		hasErr   bool
		err      error
	}{
		{
			"simple",
			&GroupExtractor{
				Name:     "product",
				Selector: ".product--full",
				Required: true,
				Extractors: map[string]*Extractor{
					"title": {
						FieldName:    "title",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__title",
					},
					"price": {
						FieldName:    "price",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__price--amount",
					},
				},
			},
			[]map[string]any{
				{
					"title": []any{"Main Product Title"},
					"price": []any{"99.99"},
				},
				{
					"title": []any{"Another Product Title"},
					"price": []any{"49.99"},
				},
				{
					"title": []any{"Third Product Title"},
					"price": []any{"0.99"},
				},
			},
			false,
			nil,
		},
		{
			"required field are empty",
			&GroupExtractor{
				Name:     "product",
				Selector: ".product--not-exists",
				Required: true,
				Extractors: map[string]*Extractor{
					"title": {
						FieldName:    "title",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__title",
					},
					"price": {
						FieldName:    "price",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__price--amount",
					},
				},
			},
			nil,
			true,
			ErrRequiredFieldMissing,
		},
	}

	// use testify assert
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {

			// check error
			skipExpectedErr := func(actual error) bool {

				if actual == nil {
					// continue test case execution
					return false
				}

				// force error if not expected
				if !c.hasErr {
					assert.NoError(t, actual)
				}

				// check error instance
				if c.err != nil {
					assert.ErrorIs(t, c.err, actual)
				}

				// stops test case execution
				return true
			}

			fn, err := Build(c.group)
			if skip := skipExpectedErr(err); skip {
				return
			}

			// compare values
			read := strings.NewReader(GetTestFieldsHTML(t))
			dom, err := goquery.NewDocumentFromReader(read)
			require.NoError(t, err)

			values, err := fn(dom.Selection)
			if skip := skipExpectedErr(err); skip {
				return
			}

			assert.Equal(t, c.expected, values)
		})
	}
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

func GetTestFieldsHTML(t *testing.T) string {

	t.Helper()

	// open file `article_test.html` return as string
	f, err := os.Open("fields_test.html")
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
}

func TestEntityTransformNewDocumentFromReaderError(t *testing.T) {
	extractor := &Extractor{
		InputFormat:  "html",
		OutputFormat: []string{"text"},
	}
	input := "<div>Hello  world!</div>"
	output := EntryTransform(extractor, input)
	assert.Equal(t, "Hello world!", output)
}

func (e *Extractor) Map() map[string]any {
	return map[string]any{
		"FieldName":    e.FieldName,
		"Limit":        e.Limit,
		"Required":     e.Required,
		"InputFormat":  e.InputFormat,
		"OutputFormat": e.OutputFormat,
		"Selector":     e.Selector,
		"BetweenStart": e.BetweenStart,
		"BetweenEnd":   e.BetweenEnd,
		"FinalRegex":   e.FinalRegex,
		"Multiline":    e.Multiline,
	}
}

func TestExtractorMap(t *testing.T) {

	e := &Extractor{
		FieldName:    "title",
		Limit:        1,
		InputFormat:  "html",
		OutputFormat: []string{"text"},
		Selector:     ".product__title",
	}

	m := e.Map()

	assert.Equal(t, "title", m["FieldName"])
	assert.Equal(t, 1, m["Limit"])
	assert.Equal(t, "html", m["InputFormat"])
	assert.Equal(t, []string{"text"}, m["OutputFormat"])
	assert.Equal(t, ".product__title", m["Selector"])
}

func (e *GroupExtractor) Map() map[string]any {
	return map[string]any{
		"Name":       e.Name,
		"Selector":   e.Selector,
		"Required":   e.Required,
		"Extractors": e.Extractors,
	}
}

func TestGroupExtractorMap(t *testing.T) {

	e := &GroupExtractor{
		Name:     "product",
		Selector: ".product--full",
		Required: true,
		Extractors: map[string]*Extractor{
			"title": {
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__title",
			},
			"price": {
				FieldName:    "price",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__price--amount",
			},
		},
	}

	m := e.Map()

	assert.Equal(t, "product", m["Name"])
	assert.Equal(t, ".product--full", m["Selector"])
	assert.True(t, m["Required"].(bool))
	assert.Len(t, m["Extractors"].(map[string]*Extractor), 2)

	title := m["Extractors"].(map[string]*Extractor)["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := m["Extractors"].(map[string]*Extractor)["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

//
//<!DOCTYPE html>
//<html lang="en">
//<head>
//<meta charset="UTF-8">
//<meta name="viewport" content="width=device-width, initial-scale=1.0">
//<title>Product Page Example</title>
//<-- Schema.org meta data -->
//<meta itemprop="name" content="Main Product Title">
//<meta itemprop="description" content="This is an example product description. It provides details about the product, its features, and benefits.">
//<meta itemprop="image" content="product-image.jpg">
//</head>
//<body>
//
//<!-- Single Product Section -->
//<div class="product product--full product-123">
//<h1 class="product__title">Main Product Title</h1>
//<div class="product__details">
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">99.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__category">
//Category: Magic wands
//</p>
//<p class="product__description">This is an example product description. It provides details about the product, its features, and benefits.</p>
//<ul class="product__features">
//<li class="product__feature">Feature 1</li>
//<li class="product__feature">Feature 2</li>
//<li class="product__feature">Feature 3</li>
//</ul>
//<div class="product__rating">
//<span class="product__rating-stars">★★★★☆</span>
//<span class="product__rating-count">(25 reviews)</span>
//</div>
//<button class="product__cart">Add to Cart</button>
//</div>
//</div>
//
//<!-- Multiple Products Section -->
//<div class="products">
//<div class="product product--related product-124">
//<h2 class="product__title">Another Product Title</h2>
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">49.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__description">Another product description providing essential details.</p>
//<button class="product__cart">Add to Cart</button>
//</div>
//
//<div class="product product--related product-125">
//<h2 class="product__title">Third Product Title</h2>
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">0.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__description">A brief description of the third product.</p>
//<button class="product__cart">Add to Cart</button>
//</div>
//</div>
//
//<!-- Additional Information Section -->
//<div class="additional-info">
//<h2 class="additional-info__title">Additional Information</h2>
//<p class="additional-info__content">This section contains additional information related to the products, such as shipping details, return policy, and FAQs.</p>
//</div>
//</body>
//</html>

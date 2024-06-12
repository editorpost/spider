package fields

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"strings"
)

// EntriesAsString extracts text or HTML content from a goquery.Selection based on the Field configuration.
// It processes the selection using the specified input format and optional CSS selector.
//
// Parameters:
//   - f (*Field): A pointer to an Field struct containing the extraction configuration.
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
//	extractor := &Field{
//	    InputFormat: "text",
//	    Selector:    "p",
//	}
//	results := EntriesAsString(extractor, doc.Selection)
//	fmt.Println(results) // Output: ["Hello", "world!"]
func EntriesAsString(f *Field, sel *goquery.Selection) []string {

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
// It processes each string according to the Field configuration.
//
// Parameters:
//   - f (*Field): A pointer to an Field struct containing the transformation configuration.
//   - entries ([]string): A slice of strings to be transformed.
//
// Returns:
//   - []string: A slice of transformed strings after applying the EntryTransform function to each entry.
//
// Example:
//
//	extractor := &Field{
//	    InputFormat:  "html",
//	    OutputFormat: []string{"text"},
//	}
//	inputs := []string{"<div>Hello  world!</div>", "<p>Go is  awesome!</p>"}
//	outputs := EntriesTransform(extractor, inputs)
//	fmt.Println(outputs) // Output: ["Hello world!", "Go is awesome!"]
func EntriesTransform(f *Field, entries []string) []string {

	for i := range entries {
		entries[i] = EntryTransform(f, entries[i])
	}

	return entries
}

// EntryTransform transforms the input value based on the Field configuration.
// It processes the value according to the specified input and output formats,
// and applies string manipulation to clean up the output.
//
// Parameters:
//   - f (*Field): A pointer to an Field struct containing the transformation configuration.
//   - value (string): The input string to be transformed.
//
// Returns:
//   - string: The transformed string after applying the specified input and output format transformations and cleaning.
//
// Example:
//
//	extractor := &Field{
//	    InputFormat:  "html",
//	    OutputFormat: []string{"text"},
//	}
//	input := "<div>Hello  world!</div>"
//	output := EntryTransform(extractor, input)
//	fmt.Println(output) // Output: "Hello world!"
func EntryTransform(f *Field, value string) string {

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
				// and fieldExtract text
				value = fromHTML.Text()
			}

		case "html":
			// do nothing, value is already html or text
		}
	}

	return ReduceSpaces(value)
}

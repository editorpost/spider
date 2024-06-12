package fields

import (
	"errors"
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

	// Selector is a css selector to find the element or limit area for between/regex.
	// optional
	Selector string `json:"Selector"`

	// between is a pair of strings to find the element.
	// In case if Selector is not provided, between applied on whole codfield.

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

	between *regexp.Regexp
	final   *regexp.Regexp

	// Children is a map of sub-field names to their corresponding Extractor configurations.
	// required
	Children []*Extractor `json:"Children" validate:"optional,dive"`
}

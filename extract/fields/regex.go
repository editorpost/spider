package fields

import (
	"github.com/editorpost/donq/pkg/valid"
	"github.com/samber/lo"
	"regexp"
)

// RegexPipes filters and processes selections through multiple regular expressions.
// It extracts substrings from each selection that match any of the given regular expressions.
//
// Parameters:
//   - selections ([]string): A slice of strings to be processed.
//   - expressions (...*regexp.Regexp): A variadic list of regular expressions to apply to the selections.
//
// Returns:
//   - []string: A slice of strings containing all the matches from the selections,
//     processed through the provided regular expressions.
//
// Example:
//
//	selections := []string{"Item 123", "Product 456", "Service 789"}
//	re1 := regexp.MustCompile(`\d+`)
//	re2 := regexp.MustCompile(`[A-Za-z]+`)
//	result := RegexPipes(selections, re1, re2)
//	fmt.Println(result) // Output: ["123", "Item", "456", "Product", "789", "Service"]
func RegexPipes(selections []string, expressions ...*regexp.Regexp) (entries []string) {

	expressions = lo.Filter(expressions, func(v *regexp.Regexp, i int) bool {
		return v != nil
	})

	if len(expressions) == 0 {
		return selections
	}

	for _, selection := range selections {
		for _, expr := range expressions {
			entries = append(entries, RegexExtract(expr, selection)...)
		}
	}

	return entries
}

// RegexCompile compiles regular expressions based on the Extractor configuration.
// It validates the Extractor and then compiles regex patterns for Between and Final matching.
//
// Parameters:
//   - f (*Extractor): A pointer to an Extractor struct containing the configuration for regex compilation.
//
// Returns:
//   - Between (*regexp.Regexp): The compiled regular expression for matching text Between BetweenStart and BetweenEnd.
//   - regex (*regexp.Regexp): The compiled user-defined regular expression from FinalRegex.
//   - err (error): An error that occurred during validation or regex compilation, or nil if successful.
//
// Example:
//
//	extractor := &Extractor{
//	    BetweenStart: "start",
//	    BetweenEnd:   "end",
//	    FinalRegex:   `\d+`,
//	}
//	Between, regex, err := RegexCompile(extractor)
//	if err != nil {
//	    log.Fatalf("Failed to compile regex: %v", err)
//	}
//	fmt.Println(Between) // Output: regexp that matches "start" and "end" with any text Between
//	fmt.Println(regex)   // Output: regexp that matches one or more digits
func RegexCompile(f *Extractor) (between *regexp.Regexp, regex *regexp.Regexp, err error) {

	// validate the fExtractor
	if err = valid.Struct(f); err != nil {
		return
	}

	// regex retrieves the text Between BetweenStart and next entry BetweenEnd
	if f.BetweenStart != "" && f.BetweenEnd != "" {
		between = regexp.MustCompile(regexp.QuoteMeta(f.BetweenStart) + "(?s)(.*?)" + regexp.QuoteMeta(f.BetweenEnd))
	}

	// compile the user defined regex, if provided
	if f.FinalRegex != "" {
		if regex, err = regexp.Compile(f.FinalRegex); err != nil {
			return
		}
	}

	return
}

// RegexExtract extracts substrings from the input data that match the given regular expression.
// It returns a slice of strings containing the first capturing group from each match.
//
// Parameters:
//   - re (*regexp.Regexp): The compiled regular expression to apply to the input data.
//   - data (string): The input string to be searched.
//
// Returns:
//   - []string: A slice of strings containing the first capturing group from each match,
//     or nil if the input data is an empty string.
//
// Example:
//
//	re := regexp.MustCompile(`(\d+)`)
//	data := "Item 123, Item 456"
//	result := RegexExtract(re, data)
//	fmt.Println(result) // Output: ["123", "456"]
func RegexExtract(re *regexp.Regexp, data string) []string {

	if data == "" {
		return nil
	}

	matches := re.FindAllStringSubmatch(data, -1)

	var values []string
	for _, match := range matches {
		if len(match) > 1 {
			values = append(values, match[1])
		}
	}

	return values
}

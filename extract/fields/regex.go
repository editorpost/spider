package fields

import (
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

// RegexCompile compiles regular expressions based on the Field configuration.
// It validates the Field and then compiles regex patterns for between and final matching.
//
// Parameters:
//   - f (*Field): A pointer to an Field struct containing the configuration for regex compilation.
//
// Returns:
//   - between (*regexp.Regexp): The compiled regular expression for matching text between BetweenStart and BetweenEnd.
//   - regex (*regexp.Regexp): The compiled user-defined regular expression from FinalRegex.
//   - err (error): An error that occurred during validation or regex compilation, or nil if successful.
//
// Example:
//
//	field := &Field{
//	    BetweenStart: "start",
//	    BetweenEnd:   "end",
//	    FinalRegex:   `\d+`,
//	}
//	between, regex, err := RegexCompile(field)
//	if err != nil {
//	    log.Fatalf("Failed to compile regex: %v", err)
//	}
//	fmt.Println(between) // Output: regexp that matches "start" and "end" with any text between
//	fmt.Println(regex)   // Output: regexp that matches one or more digits
func RegexCompile(f *Field) (between *regexp.Regexp, regex *regexp.Regexp, err error) {

	// regex retrieves the text between BetweenStart and next entry BetweenEnd
	if f.BetweenStart != "" && f.BetweenEnd != "" {
		between = regexp.MustCompile(regexp.QuoteMeta(f.BetweenStart) + "(?s)(.*?)" + regexp.QuoteMeta(f.BetweenEnd))
	}

	// compile the user defined regex, if provided
	if f.FinalRegex == "" {
		return
	}

	regex, err = regexp.Compile(f.FinalRegex)
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

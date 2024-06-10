package fields

import "strings"

// ReduceSpaces applies string manipulation methods to the input text.
// It removes any leading and trailing whitespace and replaces any
// occurrences of double spaces with a single space.
//
// Parameters:
//   - text (string): The input string to be processed.
//
// Returns:
//   - string: The processed string with no leading or trailing spaces and
//     all double spaces replaced with single spaces.
//
// Example:
//
//	input := "  This  is   an  example.  "
//	output := ReduceSpaces(input)
//	fmt.Println(output) // Output: "This is an example."
func ReduceSpaces(text string) string {

	text = strings.TrimSpace(text)

	// remove all double spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}

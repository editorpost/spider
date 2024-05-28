package collect

import (
	"net/url"
	"strings"
)

// PlaceholdersToRegex transforms a given URL pattern into a regular expression string for matching URLs.
// The function supports specific placeholders within the pattern to capture dynamic parts of URLs:
//
//   - {dir}: Matches any sequence of characters except slashes (/), representing a directory name in the URL path.
//   - {any}: Matches any sequence of characters, including an empty sequence, allowing for the most flexibility.
//   - {some}: Matches any non-empty sequence of characters, similar to {any} but requires at least one character.
//   - {num}: Matches any sequence of digits, useful for numeric identifiers in URLs.
//
// If the pattern does not contain any placeholders and is not empty, it is treated as a regular expression.
// This allows for advanced matching scenarios where specific regular expression features are needed.
// Patterns without placeholders are expected to be correctly escaped to match literal characters in URLs.
// An empty pattern matches nothing, which can be used to disable a particular rule or filter.
//
// The function automatically escapes certain characters (e.g., '/', '.') that could be misinterpreted by the regex engine
// when placeholders are used. It ensures that the generated regex pattern matches the entire string from start to end.
//
// Example:
//
//	pattern := "https://example.com/articles/{dir}/{any}"
//	regexStr := PlaceholdersToRegex(pattern)
//	// regexStr is now a regular expression string that can be used to match URLs following the specified pattern.
func PlaceholdersToRegex(pattern string) string {

	// no placeholders
	// hope it is a ready-made regex pattern or nothing
	if !strings.Contains(pattern, "{") {
		return pattern
	}

	str := pattern

	// Make sure to escape characters that could be interpreted by regex engine
	str = strings.ReplaceAll(str, "/", "\\/")
	str = strings.ReplaceAll(str, ".", "\\.")

	// Replace tokens with corresponding regex patterns
	str = strings.ReplaceAll(str, "{dir}", "([^/]+)")
	str = strings.ReplaceAll(str, "{any}", "(.*)")
	str = strings.ReplaceAll(str, "{some}", "(.+)")
	str = strings.ReplaceAll(str, "{num}", "(\\d+)")

	str = "^" + str + "$" // Match entire string

	return str
}

// MustHostname from url
func MustHostname(fromURL string) string {

	uri, err := url.Parse(fromURL)
	if err != nil {
		panic(err)
	}

	return uri.Hostname()
}

// MustRootUrl return the root url
// e.g. https://example.com/articles/1234/5678 => https://example.com
func MustRootUrl(fromURL string) string {

	uri, err := url.Parse(fromURL)
	if err != nil {
		panic(err)
	}

	return uri.Scheme + "://" + uri.Host // no port explicitly
}

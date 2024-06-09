package config

import (
	"net/url"
	"regexp"
	"strings"
)

// RegexPattern transforms a given URL pattern into a regular expression string for matching URLs.
// The function supports specific placeholders within the pattern to capture dynamic parts of URLs:
//
//   - {dir}: Matches any sequence of characters except slashes (/), representing a directory name in the URL path.
//   - {any}: Matches any sequence of characters, including an empty sequence, allowing for the most flexibility.
//   - {some}: Matches any non-empty sequence of characters, similar to {any} but requires at least one character.
//   - {num}: Matches any sequence of digits, useful for numeric identifiers in URLs.
//   - {one,two,three}: Matches any of the specified strings, separated by commas, allowing for multiple options.
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
//	regexStr := RegexPattern(pattern)
//	// regexStr is now a regular expression string that can be used to match URLs following the specified pattern.
func RegexPattern(pattern string) string {

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

	// find substrings with comma inside {one,two} or {one,two, three} or {one, two, three}
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		// replace {one,two,three} with (one|two|three)
		match[1] = strings.ReplaceAll(match[1], " ", "")
		pat := "(" + strings.ReplaceAll(match[1], ",", "|") + ")"
		str = strings.ReplaceAll(str, match[0], pat)
	}

	str = "^" + str + "$" // Match entire string

	return str
}

func ContentLikeURL(urlStr string) bool {
	allowedExtensions := map[string]bool{
		".php":   true,
		".xhtml": true,
		".shtml": true,
		".cfm":   true,
		".html":  true,
		".htm":   true,
		".asp":   true,
		".aspx":  true,
		".jsp":   true,
		".jspx":  true,
	}

	// Parse the URL to extract the path
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Extract the file extension if present
	path := parsedURL.Path
	if dotIndex := strings.LastIndex(path, "."); dotIndex != -1 {
		ext := path[dotIndex:]
		allowed := allowedExtensions[ext] // True if allowed, false otherwise
		return allowed
	}

	// True if no file extension is present
	return true
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

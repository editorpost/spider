
# Disclaimer. Don't use it. Under development.

# Usage as Windmill Script

Ensure you have the Windmill Mongodb resource `f/spider/resource/mongodbdb` available in your Windmill environment.

Expected arguments:
```json 
{
  "StartURL": "https://news.com/news/business/new-offshore-zones-promoted",
  "AllowedURL": "https://news.com{any}",
  "EntityURL": "https://news.com/news/{dir}/{some}",
  "EntitySelector": ".node-article--full",
  "UseBrowser": false,
  "Depth": 1
}
```

### Start task example:
```go 
package inner

import (
	"github.com/editorpost/spider"
)

func main(crawler interface{}) (interface{}, error) {
	//require github.com/editorpost/spider v0.0.7
	return 0, spider.StartWith(crawler)
}
```

Note! You can use the specific version of the library providing comment with `//require repo/pkg v0.0.7`. The version is specified for example only.
Use `git rev-parse HEAD` to get the current version of the library.

### Initialization with user parameters:
```go
package inner

import (
	"github.com/editorpost/spider"
)

func main(
	name string,
	startURL string,
	allowedURL string,
	entityURL string,
	entitySelector string,
	depth int,
	useBrowser bool,
) (interface{}, error) {

	//require github.com/editorpost/spider v0.0.7
	args := &spider.Args{
		Name:           name,
		StartURL:       startURL,
		AllowedURL:     allowedURL,
		EntityURL:      entityURL,
		EntitySelector: entitySelector,
		UseBrowser:     useBrowser,
		Depth:          depth,
		MongoDbResource: "f/spider/resource/mongodb",
	}

	return name, spider.Start(args)
}
```

## URL Pattern Placeholders

In defining URL patterns for routing, filtering, or matching purposes, our system supports the use of specific placeholders within the pattern strings. These placeholders allow for dynamic matching against various URL structures. When a URL pattern includes one or more of these placeholders, it is transformed into a regular expression that matches the corresponding URL structure.

### Available Placeholders

- `{dir}`: Matches any sequence of characters except for slashes (`/`). Used to represent directory names in a URL path.
- `{any}`: Matches any sequence of characters, including an empty sequence. This is the most flexible placeholder and can match any part of a URL.
- `{some}`: Matches any non-empty sequence of characters. Similar to `{any}`, but requires at least one character to be present.
- `{num}`: Matches any sequence of digits. Useful for matching numeric identifiers within URLs.

### Placeholder Usage

Placeholders can be inserted into URL patterns to specify the parts of the URL that can vary. For example:

- `https://example.com/articles/{dir}/{some}`: Matches URLs that follow the structure of two path segments following `/articles/`, where the first segment can be any directory name, and the second must be a non-empty sequence.
- `https://example.com/products/{num}/details`: Matches URLs that include a numeric product ID followed by `/details`.

### Patterns Without Placeholders

If a URL pattern does not contain any of the specified placeholders, the system interprets the pattern in two ways:

1. **As a Regular Expression (If Not Empty):** If the pattern is not empty and contains no placeholders, it is treated as a ready-made regular expression. This allows for advanced matching scenarios where specific regular expression features are needed. When defining such patterns, ensure they are correctly escaped and adhere to regular expression syntax rules.

2. **Literal Match (If Empty):** An empty pattern matches nothing. This can be used to effectively disable a particular matching rule or filter.

### Escaping Special Characters

When placeholders are not used, and the pattern is intended as a regular expression, special characters must be escaped to prevent them from being interpreted as regex operators. For instance, `.` should be written as `\\.` and `/` as `\\/` to match these characters literally in URLs.

This flexible system of placeholders and direct regular expression input allows for precise control over URL matching, accommodating a wide range of routing and filtering requirements.
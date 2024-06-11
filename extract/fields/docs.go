/*
Package fields provides fields.Field to extract custom entity fields.

Simple configuration for extracting data from HTML or text content:
1. **CSS Selectors with goquery**:
  - Utilize jQuery-like CSS selectors for straightforward element selection.
  - The most minimal and common case for extracting elements from HTML.

2. **High-Level Regex Abstraction**:
  - `Between` allows extraction of content between specified start and end strings.
  - Flexible for both HTML and text input formats.
  - Useful for structured data extraction where content is surrounded by predictable markers.

3. **Comprehensive Regular Expressions**:
  - Supports user-defined regular expressions for advanced extraction needs.
  - Expressions are precompiled for efficiency and reused across extractions.
  - Can be combined with `Between` for complex scenarios.

4. **Multiple Entries Extraction**:
  - Capable of extracting multiple entries from the input.
  - Handles lists, repeated elements, and other multi-entry scenarios seamlessly.

5. **Builder Limiting**:
  - Configurable limit on the number of entries to extract.
  - Zero means no limit, and entries are stored as a list.

6. **String Manipulation**:
  - Trims leading and trailing spaces.
  - Reduces multiple spaces to a single space.
  - Ensures clean and consistent output format.

7. **Transformation Configuration**:
  - InputFormat can be set to "text" or "html" for parsing and processing accordingly.
  - OutputFormat supports transformations to "text", "html", or "json".
  - Formatters are applied in the order specified in the OutputFormat list.

8. **Error Handling**:
  - Provides meaningful errors for required fields and validation issues.
  - Skips entity extraction if a required field is missing.

9. **Data Cleaning**:

  - Removes empty entries from the results.

  - Ensures uniqueness by removing duplicate entries.

    10. **Example Usage**:
    ```go
    package main

    import (
    "fmt"
    "log"
    "strings"

    "github.com/PuerkitoBio/goquery"
    )

    func main() {
    extractor := &Field{
    Name:    "example",
    InputFormat:  "html",
    OutputFormat: []string{"text"},
    Selector:     "p",
    Cardinality:        2,
    Required:     true,
    }

    extractFn, err := extractor.Builder()
    if err != nil {
    log.Fatalf("Failed to build extractor: %v", err)
    }

    htmlContent := "<div><p>Hello</p><p>world!</p></div>"
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
    if err != nil {
    log.Fatalf("Failed to parse HTML: %v", err)
    }

    results, err := extractFn(doc.Selection)
    if err != nil {
    log.Fatalf("Extraction error: %v", err)
    }

    fmt.Println(results) // Output: ["Hello", "world!"]
    }
*/
package fields

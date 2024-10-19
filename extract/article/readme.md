### `article` Package Documentation

#### Overview

The `article` package in the `editorpost/spider` repository is responsible for extracting articles from HTML content. It converts the HTML into structured data, including the article's text, title, author, summary, and publication date.

#### Principles

- **HTML Parsing**: The package uses libraries like `goquery` and `go-readability` to parse and extract information from HTML.
- **Markdown Conversion**: Extracted HTML is converted to Markdown format using the `html-to-markdown` library.
- **Normalization**: Extracted data is normalized to ensure consistency and validity.

#### Key Functions

1. **Article**: Main function that extracts the article data and sets it in the payload.
2. **ArticleSelection**: Combines head tags and article HTML selection to prepare for readability extraction.
3. **ArticleFromHTML**: Extracts article data from HTML content.
4. **ArticleSelectionToMarkup**: Converts HTML selection to Markdown.
5. **AbsoluteUrl**: Converts relative URLs to absolute URLs.
6. **readabilityArticle**: Uses `go-readability` to extract readable content from HTML.
7. **distillArticle**: Uses `go-domdistiller` to extract content from HTML.
8. **HTMLToMarkdown**: Converts HTML content to Markdown format.
9. **legacyPublished** and **legacyAuthor**: Fallback methods to extract publication date and author from legacy HTML structures.

#### Configuration Features

- **Extract Selector**: Configurable selector used to target specific parts of the HTML for extraction.
- **Proxy Settings**: Can be integrated with proxy settings to manage and rotate proxies during extraction.

#### Caveats

- **HTML Structure Dependence**: The extraction heavily depends on the structure of the HTML. Any changes in the HTML structure can affect the extraction process.
- **Error Handling**: The package logs warnings and errors during extraction but continues to attempt extraction for robustness.
- **Incomplete Data**: In case of missing or incorrect HTML elements, fallback methods are used to extract data, which may not always be accurate.

#### Error Handling and Retries

- **Error Logging**: Errors are logged using `slog.Warn` for debugging purposes.
- **Retries**: The package does not implement automatic retries but handles errors gracefully, ensuring that extraction continues as much as possible.

#### Important Notes

- **Performance**: The extraction process may be time-consuming for large HTML documents or complex structures.
- **Dependencies**: The package relies on several third-party libraries (`goquery`, `go-readability`, `domdistiller`, `html-to-markdown`) for its functionality.

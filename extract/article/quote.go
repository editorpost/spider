package article

// Quote represents a quote from social media in the article.
type Quote struct {
	// Text is the text of the quote.
	// This field is required.
	Text string `json:"text" validate:"required"`

	// Author is the author of the quote.
	// This field is required and should be between 1 and 255 characters long.
	Author string `json:"author" validate:"required,min=1,max=255"`

	// Source is the source URL of the quote.
	// This field is required and should be a valid URL.
	Source string `json:"source" validate:"required,url"`

	// Platform is the platform of the quote (e.g., Twitter, Facebook).
	// This field is required and should be between 1 and 50 characters long.
	Platform string `json:"platform" validate:"required,min=1,max=50"`
}

// Map converts the Quote struct to a map[string]any.
func (q *Quote) Map() map[string]any {
	return map[string]any{
		"text":     q.Text,
		"author":   q.Author,
		"source":   q.Source,
		"platform": q.Platform,
	}
}

// NewQuoteFromMap creates a Quote from a map[string]any, validates it, and returns a pointer to the Quote or an error.
func NewQuoteFromMap(m map[string]any) (*Quote, error) {
	quote := &Quote{
		Text:     getString(m, "text"),
		Author:   getString(m, "author"),
		Source:   getString(m, "source"),
		Platform: getString(m, "platform"),
	}

	err := validate.Struct(quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

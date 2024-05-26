package article

// Video represents a video in the article.
type Video struct {
	// URL is the URL of the video.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url"`

	// EmbedCode is the embed code for the video.
	// This field is optional.
	EmbedCode string `json:"embed_code,omitempty"`

	// Caption is the caption for the video.
	// This field is optional.
	Caption string `json:"caption,omitempty"`
}

// Map converts the Video struct to a map[string]any.
func (v *Video) Map() map[string]any {
	return map[string]any{
		"url":        v.URL,
		"embed_code": v.EmbedCode,
		"caption":    v.Caption,
	}
}

// NewVideoFromMap creates a Video from a map[string]any, validates it, and returns a pointer to the Video or an error.
func NewVideoFromMap(m map[string]any) (*Video, error) {
	video := &Video{
		URL:       getString(m, "url"),
		EmbedCode: getString(m, "embed_code"),
		Caption:   getString(m, "caption"),
	}

	err := validate.Struct(video)
	if err != nil {
		return nil, err
	}

	return video, nil
}

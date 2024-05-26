package article

// SocialProfile represents a social media profile of an author.
type SocialProfile struct {
	// Platform is the platform of the social profile (e.g., Twitter, Facebook).
	// This field is required and should be between 1 and 50 characters long.
	Platform string `json:"platform" validate:"required,min=1,max=50"`

	// URL is the URL of the social profile.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url"`
}

// Map converts the SocialProfile struct to a map[string]any.
func (sp *SocialProfile) Map() map[string]any {
	return map[string]any{
		"platform": sp.Platform,
		"url":      sp.URL,
	}
}

// NewSocialProfileFromMap creates a SocialProfile from a map[string]any, validates it, and returns a pointer to the SocialProfile or an error.
func NewSocialProfileFromMap(m map[string]any) (*SocialProfile, error) {
	profile := &SocialProfile{
		Platform: getString(m, "platform"),
		URL:      getString(m, "url"),
	}

	err := validate.Struct(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

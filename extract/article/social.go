package article

import (
	"github.com/go-playground/validator/v10"
	"log/slog"
)

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

// SocialProfile represents a social media profile of an author.
type SocialProfile struct {
	// Platform is the platform of the social profile (e.g., Twitter, Facebook).
	// This field is required and should be between 1 and 50 characters long.
	Platform string `json:"platform" validate:"max=255"`

	// URL is the URL of the social profile.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url,max=4096"`
}

// Normalize validates and trims the fields of the SocialProfile.
func (s *SocialProfile) Normalize() {

	s.Platform = TrimToMaxLen(s.Platform, 255)
	s.URL = TrimToMaxLen(s.URL, 4096)

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in SocialProfile", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
			*s = SocialProfile{}
		}
	}
}

// Map converts the SocialProfile struct to a map[string]any.
func (s *SocialProfile) Map() map[string]any {
	return map[string]any{
		"platform": s.Platform,
		"url":      s.URL,
	}
}

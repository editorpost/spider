package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

// init go validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// NewArticle creates a new Article with the provided data and returns a pointer to the Article.
func NewArticle() *Article {
	return &Article{
		ID:                   uuid.New().String(),
		Language:             "en",
		Tags:                 []string{},
		Images:               []Image{},
		Videos:               []Video{},
		Quotes:               []Quote{},
		AuthorSocialProfiles: []SocialProfile{},
	}
}

// Article represents a news article with various types of content.
// This structure provides a flexible and universal foundation for storing and working with various types of content,
// allowing for easy creation and modification of articles, as well as integration of media and social elements.
type Article struct {
	ID                   string          `json:"article__id" validate:"required,uuid4"`
	Title                string          `json:"article__title" validate:"required,max=255"`
	Byline               string          `json:"article__byline" validate:"max=255"`
	Content              string          `json:"article__content" validate:"required"`
	TextContent          string          `json:"article__text_content" validate:"required"`
	Excerpt              string          `json:"article__excerpt" validate:"max=500"`
	PublishDate          time.Time       `json:"article__publish_date" validate:"required"`
	ModifiedDate         time.Time       `json:"article__modified_date"`
	Images               []Image         `json:"article__images"`
	Videos               []Video         `json:"article__videos"`
	Quotes               []Quote         `json:"article__quotes"`
	Tags                 []string        `json:"article__tags"`
	Source               string          `json:"article__source" validate:"omitempty,url"`
	Language             string          `json:"article__language" validate:"required,len=2"`
	Category             string          `json:"article__category" validate:"max=255"`
	SiteName             string          `json:"article__site_name" validate:"max=255"`
	AuthorSocialProfiles []SocialProfile `json:"article__author_social_profiles"`
}

// Normalize validates the Article and its nested structures, logs any validation errors, and clears invalid fields.
func (a *Article) Normalize() {
	validate := validator.New()

	err := validate.Struct(a)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			slog.Info("Validation error", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))

			// Clear invalid fields
			switch err.Namespace() {
			case "Article.ID":
				a.ID = ""
			case "Article.Title":
				a.Title = ""
			case "Article.Byline":
				a.Byline = ""
			case "Article.Content":
				a.Content = ""
			case "Article.TextContent":
				a.TextContent = ""
			case "Article.Excerpt":
				a.Excerpt = ""
			case "Article.PublishDate":
				a.PublishDate = time.Time{}
			case "Article.ModifiedDate":
				a.ModifiedDate = time.Time{}
			case "Article.Source":
				a.Source = ""
			case "Article.Language":
				a.Language = ""
			case "Article.Category":
				a.Category = ""
			case "Article.SiteName":
				a.SiteName = ""
			}
		}
	}

	// Validate nested structures
	for i, image := range a.Images {
		err := validate.Struct(image)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				slog.Info("Validation error in Image", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
				a.Images[i] = Image{}
			}
		}
	}

	for i, video := range a.Videos {
		err := validate.Struct(video)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				slog.Info("Validation error in Video", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
				a.Videos[i] = Video{}
			}
		}
	}

	for i, quote := range a.Quotes {
		err := validate.Struct(quote)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				slog.Info("Validation error in Quote", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
				a.Quotes[i] = Quote{}
			}
		}
	}

	for i, profile := range a.AuthorSocialProfiles {
		err := validate.Struct(profile)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				slog.Info("Validation error in SocialProfile", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
				a.AuthorSocialProfiles[i] = SocialProfile{}
			}
		}
	}
}

// Map converts the Article struct to a map[string]any, including nested structures.
func (a *Article) Map() map[string]any {
	images := make([]map[string]any, len(a.Images))
	for i, image := range a.Images {
		images[i] = image.Map()
	}

	videos := make([]map[string]any, len(a.Videos))
	for i, video := range a.Videos {
		videos[i] = video.Map()
	}

	quotes := make([]map[string]any, len(a.Quotes))
	for i, quote := range a.Quotes {
		quotes[i] = quote.Map()
	}

	socialProfiles := make([]map[string]any, len(a.AuthorSocialProfiles))
	for i, profile := range a.AuthorSocialProfiles {
		socialProfiles[i] = profile.Map()
	}

	return map[string]any{
		"article__id":                     a.ID,
		"article__title":                  a.Title,
		"article__byline":                 a.Byline,
		"article__content":                a.Content,
		"article__text_content":           a.TextContent,
		"article__excerpt":                a.Excerpt,
		"article__images":                 images,
		"article__videos":                 videos,
		"article__quotes":                 quotes,
		"article__publish_date":           a.PublishDate,
		"article__modified_date":          a.ModifiedDate,
		"article__tags":                   a.Tags,
		"article__source":                 a.Source,
		"article__language":               a.Language,
		"article__category":               a.Category,
		"article__site_name":              a.SiteName,
		"article__author_social_profiles": socialProfiles,
	}
}

// NewArticleFromMap creates an Article from a map[string]any, validates it, and returns a pointer to the Article or an error.
func NewArticleFromMap(m map[string]any) (*Article, error) {
	images := make([]Image, 0)
	if imgMaps, ok := m["article__images"].([]map[string]any); ok {
		for _, imgMap := range imgMaps {
			if img, err := NewImageFromMap(imgMap); err == nil {
				images = append(images, *img)
			}
		}
	}

	videos := make([]Video, 0)
	if vidMaps, ok := m["article__videos"].([]map[string]any); ok {
		for _, vidMap := range vidMaps {
			if vid, err := NewVideoFromMap(vidMap); err == nil {
				videos = append(videos, *vid)
			}
		}
	}

	quotes := make([]Quote, 0)
	if quoteMaps, ok := m["article__quotes"].([]map[string]any); ok {
		for _, quoteMap := range quoteMaps {
			if quote, err := NewQuoteFromMap(quoteMap); err == nil {
				quotes = append(quotes, *quote)
			}
		}
	}

	socialProfiles := make([]SocialProfile, 0)
	if profileMaps, ok := m["article__author_social_profiles"].([]map[string]any); ok {
		for _, profileMap := range profileMaps {
			if profile, err := NewSocialProfileFromMap(profileMap); err == nil {
				socialProfiles = append(socialProfiles, *profile)
			}
		}
	}

	publishDate, _ := m["article__publish_date"].(time.Time)
	modifiedDate, _ := m["article__modified_date"].(time.Time)

	article := &Article{
		ID:                   getString(m, "article__id"),
		Title:                getString(m, "article__title"),
		Byline:               getString(m, "article__byline"),
		Content:              getString(m, "article__content"),
		TextContent:          getString(m, "article__text_content"),
		Excerpt:              getString(m, "article__excerpt"),
		Images:               images,
		Videos:               videos,
		Quotes:               quotes,
		PublishDate:          publishDate,
		ModifiedDate:         modifiedDate,
		Tags:                 GetStringSlice(m, "article__tags"),
		Source:               getString(m, "article__source"),
		Language:             getString(m, "article__language"),
		Category:             getString(m, "article__category"),
		SiteName:             getString(m, "article__site_name"),
		AuthorSocialProfiles: socialProfiles,
	}

	err := validate.Struct(article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetStringSlice safely extracts a slice of strings from the map or returns a zero value.
func GetStringSlice(m map[string]any, key string) []string {
	if value, ok := m[key]; ok {
		if slice, ok := value.([]string); ok {
			return slice
		}
	}
	return []string{}
}

package article_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/spider/extract/article" // Замените на фактический путь к вашему пакету
)

func init() {
	gofakeit.Seed(0)
}

func TestMinimalInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Content = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.PublishDate = time.Now()

	got, err := article.NewArticleFromMap(expected.Map())
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestFullInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Content = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.PublishDate = time.Now()
	expected.ModifiedDate = time.Now()

	// Optional fields
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)
	expected.Images = []article.Image{
		{
			URL:     gofakeit.URL(),
			AltText: gofakeit.Sentence(5),
			Width:   gofakeit.Number(800, 1920),
			Height:  gofakeit.Number(600, 1080),
			Caption: gofakeit.Sentence(10),
		},
	}
	expected.Videos = []article.Video{
		{
			URL:       gofakeit.URL(),
			EmbedCode: "<iframe src='" + gofakeit.URL() + "'></iframe>",
			Caption:   gofakeit.Sentence(10),
		},
	}
	expected.Quotes = []article.Quote{
		{
			Text:     gofakeit.Sentence(15),
			Author:   gofakeit.Name(),
			Source:   gofakeit.URL(),
			Platform: "Twitter",
		},
	}
	expected.Tags = []string{"travel", "Phuket", "Thailand"}
	expected.Source = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"
	expected.AuthorSocialProfiles = []article.SocialProfile{
		{
			Platform: "Twitter",
			URL:      gofakeit.URL(),
		},
	}

	got, err := article.NewArticleFromMap(expected.Map())
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestInvalidNestedStructureArticle(t *testing.T) {
	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Content = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.PublishDate = time.Now()
	expected.ModifiedDate = time.Now()

	// Optional fields with invalid nested structure
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)
	expected.Images = []article.Image{
		{
			URL:     "invalid-url",
			AltText: gofakeit.Sentence(5),
			Width:   gofakeit.Number(800, 1920),
			Height:  gofakeit.Number(600, 1080),
			Caption: gofakeit.Sentence(10),
		},
	}
	expected.Tags = []string{"travel", "Phuket", "Thailand"}
	expected.Source = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"

	// Convert expected Article to map and then back to Article to simulate input processing
	inputMap := expected.Map()

	// Expect the images to be nil due to invalid URL in the nested structure
	expected.Images = []article.Image{}

	got, err := article.NewArticleFromMap(inputMap)
	require.NoError(t, err)

	// To compare PublishDate and ModifiedDate separately due to possible time differences
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Byline, got.Byline)
	assert.Equal(t, expected.Content, got.Content)
	assert.Equal(t, expected.TextContent, got.TextContent)
	assert.Equal(t, expected.Excerpt, got.Excerpt)
	assert.Equal(t, expected.Images, got.Images)
	assert.WithinDuration(t, expected.PublishDate, got.PublishDate, time.Second)
	assert.WithinDuration(t, expected.ModifiedDate, got.ModifiedDate, time.Second)
	assert.Equal(t, expected.Tags, got.Tags)
	assert.Equal(t, expected.Source, got.Source)
	assert.Equal(t, expected.Language, got.Language)
	assert.Equal(t, expected.Category, got.Category)
	assert.Equal(t, expected.SiteName, got.SiteName)
}

func TestMissingRequiredFieldsArticle(t *testing.T) {
	art, err := article.NewArticleFromMap(article.NewArticle().Map())
	require.Error(t, err)
	assert.Nil(t, art)
}

func TestArticleNormalize(t *testing.T) {
	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Content = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.PublishDate = time.Now()
	expected.ModifiedDate = time.Now()

	// Optional fields with some invalid data
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)
	expected.Images = []article.Image{
		{
			URL:     "invalid-url",
			AltText: gofakeit.Sentence(5),
			Width:   gofakeit.Number(800, 1920),
			Height:  gofakeit.Number(600, 1080),
			Caption: gofakeit.Sentence(10),
		},
	}
	expected.Videos = []article.Video{
		{
			URL:       "invalid-url",
			EmbedCode: "<iframe src='invalid-url'></iframe>",
			Caption:   gofakeit.Sentence(10),
		},
	}
	expected.Quotes = []article.Quote{
		{
			Text:     "",
			Author:   gofakeit.Name(),
			Source:   "invalid-url",
			Platform: "Twitter",
		},
	}
	expected.Tags = []string{"travel", "Phuket", "Thailand"}
	expected.Source = "invalid-url"
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"
	expected.AuthorSocialProfiles = []article.SocialProfile{
		{
			Platform: "Twitter",
			URL:      "invalid-url",
		},
	}

	expected.Normalize()

	// Verify that invalid fields are cleared
	assert.Equal(t, "", expected.Images[0].URL)
	assert.Equal(t, "", expected.Videos[0].URL)
	assert.Equal(t, "", expected.Quotes[0].Text)
	assert.Equal(t, "", expected.Quotes[0].Source)
	assert.Equal(t, "", expected.Source)
	assert.Equal(t, "", expected.AuthorSocialProfiles[0].URL)
}

func TestArticleNormalizeFieldClearing(t *testing.T) {

	invalid := article.NewArticle()

	// Set required fields with valid data
	invalid.Title = gofakeit.Sentence(3)
	invalid.Content = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.PublishDate = time.Now()
	invalid.ModifiedDate = time.Now()

	// Set invalid data for optional fields
	invalid.ID = "invalid-uuid"
	invalid.Byline = gofakeit.Name()
	invalid.Excerpt = gofakeit.Sentence(10)
	invalid.Source = "invalid-url"
	invalid.Language = "inglese" // should be a valid ISO 639-1 language code
	invalid.Category = gofakeit.Sentence(2)
	invalid.SiteName = gofakeit.Sentence(2)

	valid := *invalid
	(&valid).Normalize()

	// Verify that invalid fields are cleared
	assert.Equal(t, "", valid.ID)
	assert.Equal(t, invalid.Byline, valid.Byline)   // should not be cleared since it's not required
	assert.Equal(t, invalid.Excerpt, valid.Excerpt) // should not be cleared since it's not required
	assert.Equal(t, "", valid.Source)
	assert.Equal(t, "inglese", valid.Language)
	assert.Equal(t, invalid.Category, valid.Category)
	assert.Equal(t, invalid.SiteName, valid.SiteName)
}

// TestGetStringSlice tests the GetStringSlice function in case of empty map and missing key:
func TestGetStringSlice(t *testing.T) {
	m := map[string]interface{}{}
	key := "key"
	assert.Equal(t, []string{}, article.GetStringSlice(m, key))
}

func TestTrimToMaxLen(t *testing.T) {
	s := "This is a test string with more than twenty characters."
	trimmed := article.TrimToMaxLen(s, 20)
	assert.Equal(t, "This is a test strin", trimmed)

	s = "Short string"
	trimmed = article.TrimToMaxLen(s, 20)
	assert.Equal(t, s, trimmed)
}

//Напиши `readme.md` o пакете `Article` для разработчиков и пользователей.
//
//Из документа должны быть понятны принятые архитектурные решения, ограничения и их причины. Стоит отдельно рассмотреть лимиты наложенные валидацией, так чтобы это было удобно и для проектировщика БД и для редактора материалов.
//
//Покажи минимальный инвариант в виде json.
//
//Подчеркни, что рекомендуется использовать `article.NewArticle`, чтобы иметь структуру близкую к инварианту. Перечили поля, которые необходимо добавить для достижения минимального инварианта.
//
// Создай сводку для "тех, кто спешит". Добавь оглавление со ссылками.

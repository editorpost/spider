package article_test

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/tester"
	"github.com/go-shiori/go-readability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestArticleFromPayload(t *testing.T) {

	uri := gofakeit.URL()
	filename := "../../tester/fixtures/article.html"
	payload := tester.TestPayloadWithURI(t, filename, uri)

	a, err := article.ArticleFromPayload(payload)
	require.NoError(t, err)

	// custom fallback with css selector
	assert.Equal(t, "2024-03-13", a.Published.Format("2006-01-02"))
	assert.Equal(t, "John Doe", a.Author)
	assert.Equal(t, "Пхукет в стиле вашего отдыха", a.Title)
	assert.Equal(t, uri, a.SourceURL)
	assert.Equal(t, "", a.SourceName)

	// validation
	assert.NoError(t, a.Normalize())
}

func TestArticleFromPayload_Title(t *testing.T) {

	payload := tester.TestPayload(t, "../../tester/fixtures/cases/must_article_title.html")
	a, err := article.ArticleFromPayload(payload)
	require.NoError(t, err)

	assert.Equal(t, "Рамбутан", a.Title)
}

func TestReadability(t *testing.T) {

	markup := GetArticleHTML(t)
	read, err := readability.FromReader(strings.NewReader(markup), GetArticleURL(t))
	require.NoError(t, err)

	fmt.Printf("Title   : %s\n", read.Title)
	fmt.Printf("Date    : %s\n", read.PublishedTime)
	fmt.Printf("Author  : %s\n", read.Byline)
	fmt.Printf("Length  : %d\n", read.Length)
	fmt.Printf("Genre : %s\n", read.Excerpt)
	fmt.Printf("SourceName: %s\n", read.SiteName)
	fmt.Printf("Image   : %s\n", read.Image)
	fmt.Printf("Favicon : %s\n", read.Favicon)
	fmt.Printf("ContentText:\n %s \n", read.TextContent)
	fmt.Printf("Content:\n %s \n", read.Content)
}

//goland:noinspection GoUnhandledErrorResult
func TestMarkdown(t *testing.T) {

	markup := GetArticleHTML(t)
	u := &url.URL{
		Scheme: "https",
		Host:   "thailand-news.ru",
		Path:   "/news/puteshestviya/pkhuket-v-stile-vashego-otdykha/",
	}
	read, err := readability.FromReader(strings.NewReader(markup), u)
	require.NoError(t, err)

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(read.Content)
	if err != nil {
		log.Fatal(err)
	}

	// save to file `extract_test.md`
	f, err := os.Create("article_test.md")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(markdown)
	require.NoError(t, err)
}

func GetArticleHTML(t *testing.T) string {

	t.Helper()

	// open file `article_test.html` return as string
	f, err := os.Open("article_test.html")
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
}

func GetArticleURL(t *testing.T) *url.URL {

	t.Helper()
	return &url.URL{
		Scheme: "https",
		Host:   "thailand-news.ru",
		Path:   "/news/puteshestviya/pkhuket-v-stile-vashego-otdykha/",
	}
}

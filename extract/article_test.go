package extract_test

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/editorpost/spider/extract"
	"github.com/go-shiori/dom"
	"github.com/go-shiori/go-readability"
	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestFromHTML(t *testing.T) {

	a, err := extract.ArticleFromHTML(GetArticleHTML(t), GetArticleURL(t))
	require.NoError(t, err)

	// custom fallback with css selector
	assert.Equal(t, "2024-03-13", a.Published.Format("2006-01-02"))
	assert.Equal(t, "Egor Denisov", a.Author)
	assert.Equal(t, "Пхукет в стиле вашего отдыха", a.Title)
	assert.Equal(t, "https://thailand-news.ru/news/puteshestviya/pkhuket-v-stile-vashego-otdykha/", a.SourceURL)
	assert.Equal(t, "", a.SourceName)

	// validation
	assert.NoError(t, a.Normalize())

	m := a.Map()
	_ = m
	fmt.Println(m)
}

func TestDistiller(t *testing.T) {

	markup := GetArticleHTML(t)

	// Run distiller
	result, err := distiller.ApplyForReader(strings.NewReader(markup), &distiller.Options{
		OriginalURL: GetArticleURL(t),
	})
	if err != nil {
		panic(err)
	}

	rawHTML := dom.OuterHTML(result.Node)
	fmt.Println(rawHTML)
}

func TestReadability(t *testing.T) {

	markup := GetArticleHTML(t)
	article, err := readability.FromReader(strings.NewReader(markup), GetArticleURL(t))
	require.NoError(t, err)

	fmt.Printf("Title   : %s\n", article.Title)
	fmt.Printf("Date    : %s\n", article.PublishedTime)
	fmt.Printf("Author  : %s\n", article.Byline)
	fmt.Printf("Length  : %d\n", article.Length)
	fmt.Printf("Genre : %s\n", article.Excerpt)
	fmt.Printf("SourceName: %s\n", article.SiteName)
	fmt.Printf("Image   : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Printf("ContentText:\n %s \n", article.TextContent)
	fmt.Printf("Content:\n %s \n", article.Content)
}

//goland:noinspection GoUnhandledErrorResult
func TestMarkdown(t *testing.T) {

	markup := GetArticleHTML(t)
	u := &url.URL{
		Scheme: "https",
		Host:   "thailand-news.ru",
		Path:   "/news/puteshestviya/pkhuket-v-stile-vashego-otdykha/",
	}
	article, err := readability.FromReader(strings.NewReader(markup), u)
	require.NoError(t, err)

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(article.Content)
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

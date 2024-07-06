package article_test

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/editorpost/spider/extract/article"
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

	a, err := article.ArticleFromHTML(GetArticleHTML(t), GetArticleURL(t))
	require.NoError(t, err)

	// custom fallback with css selector
	assert.Equal(t, "2024-03-13", a.Published.Format("2006-01-02"))
	assert.Equal(t, "John Doe", a.Author)
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

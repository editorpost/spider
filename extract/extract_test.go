package extract_test

import (
	"bytes"
	"fmt"
	readability "github.com/go-shiori/go-readability"
	"github.com/stretchr/testify/require"
	"net/url"
	"os"
	"testing"
)

func TestReadability(t *testing.T) {

	b, err := os.ReadFile("extract_test.html")
	require.NoError(t, err)

	pageURL, err := url.Parse("https://example.com")
	require.NoError(t, err)

	article, err := readability.FromReader(bytes.NewReader(b), pageURL)
	require.NoError(t, err)

	fmt.Printf("URL     : %s\n", pageURL)
	fmt.Printf("Title   : %s\n", article.Title)
	fmt.Printf("Author  : %s\n", article.Byline)
	fmt.Printf("Length  : %d\n", article.Length)
	fmt.Printf("Excerpt : %s\n", article.Excerpt)
	fmt.Printf("SiteName: %s\n", article.SiteName)
	fmt.Printf("Image   : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Printf("ContentText:\n %s \n", article.TextContent)
	fmt.Printf("Content:\n %s \n", article.Content)
}

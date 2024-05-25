package extract_test

import (
	"github.com/editorpost/spider/extract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImagesURLs(t *testing.T) {
	htmlStr, err := extract.LoadHTMLFromFile("article_test.html")
	require.NoError(t, err, "failed to load HTML from file")

	payload, err := extract.CreatePayload(htmlStr, "https://thailand-news.ru")
	require.NoError(t, err, "failed to create payload")

	err = extract.Article(payload)
	require.NoError(t, err, "Article function returned an error")

	err = extract.ImagesURLs(payload, extract.MinSizeFilter(300, 300))
	require.NoError(t, err, "ImagesURLs function returned an error")

	expectedImages := []string{
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Phuket.png",
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Rosewood%201.jpg",
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Nai%20Harn%202%201187x789.jpg",
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Anantara%203.jpg",
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Kata%20Rocks%204.jpg",
		"https://thailand-news.ru/sites/default/files/storage/images/2024-11/Trisara%205.jpg",
	}

	assert.ElementsMatch(t, expectedImages, payload.Data["entity__images"].([]string))
}

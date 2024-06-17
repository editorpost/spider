package collect_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCollect(t *testing.T) {

	srv := ServeFile(t, "crawler_test.html")
	defer srv.Close()

	dispatched := false

	crawler, err := collect.NewCrawler(
		&config.Args{
			StartURL:        srv.URL,
			AllowedURL:      ".*",
			Depth:           1,
			ExtractSelector: ".article--ssr",
		},
		&config.Deps{
			Extractor: func(*colly.HTMLElement, *goquery.Selection) error {
				dispatched = true
				return nil
			},
		},
	)
	require.NoError(t, err)

	err = crawler.Run()

	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestJSCollect(t *testing.T) {

	srv := ServeFile(t, "crawler_test.html")
	defer srv.Close()

	dispatched := false

	crawler, err := collect.NewCrawler(
		&config.Args{
			StartURL:        srv.URL,
			AllowedURL:      ".*",
			Depth:           1,
			ExtractSelector: ".article--js",
			UseBrowser:      true,
		},
		&config.Deps{
			Extractor: func(*colly.HTMLElement, *goquery.Selection) error {
				dispatched = true
				return nil
			},
		},
	)

	require.NoError(t, err)
	require.NoError(t, crawler.Run())
	assert.True(t, dispatched)
}

func TestMongoConfig(t *testing.T) {
	// Test the mongodb config
	validResource := map[string]interface{}{
		"db": "exampleDb",
		"servers": []interface{}{
			map[string]interface{}{
				"host": "exampleHost",
				"port": 1234.0,
			},
		},
		"credential": map[string]interface{}{
			"password": "examplePassword",
			"username": "exampleUsername",
		},
	}

	conf, err := mongodb.ConfigFromResource(validResource)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, conf, "Expected non-nil config")
}

func TestServeFile(t *testing.T) {

	srv := ServeFile(t, "crawler_test.html")
	defer srv.Close()

	// create a new request
	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	// create http client
	client := srv.Client()

	// send the request
	resp, err := client.Do(req)
	require.NoError(t, err)

	// check the response
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// read the response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body := string(b)

	// check the response body
	require.NotNil(t, body)

	// check html contains string "Hello, World!"
	require.Contains(t, body, "Hello, World!")
}

func Parse[T any](from any, to *T) error {

	m, ok := from.(map[string]interface{})
	if !ok {
		return errors.New("invalid input arguments")
	}

	// Convert the map to JSON
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal input arguments: %w", err)
	}

	// Convert the JSON to a struct
	if err = json.Unmarshal(data, to); err != nil {
		return fmt.Errorf("failed to unmarshal input arguments: %w", err)
	}

	return nil
}

// ServeFile serves the file at the given path and returns the URL
func ServeFile(t *testing.T, path string) *httptest.Server {

	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// create http server and serve the file

	srv, err := NewServer(b)
	if err != nil {
		t.Fatal(err)
	}

	// return the server URL
	return srv
}

// NewServer creates a new server that serves the given content
func NewServer(content []byte) (*httptest.Server, error) {

	// create a new server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	return srv, nil
}

package spider_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
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

const (
	testDbName = "test_test"
)

var (
	mongoCfg  *mongodb.Config
	collector *store.CollectStore
)

func TestMain(m *testing.M) {

	res := map[string]any{
		"db": "spider_meta",
		"servers": []any{
			map[string]any{
				"host": "localhost",
				"port": 27018,
			},
		},
		"credential": map[string]any{
			"username": "root",
			"password": "nopass",
		},
	}

	var err error
	if mongoCfg, err = mongodb.ConfigFromResource(res); err != nil {
		log.Fatal(err)
	}

	collector, err = store.NewCollectStore(testDbName, mongoCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestCollect(t *testing.T) {

	srv := ServeFile(t, "run_test.html")
	defer srv.Close()

	dispatched := false

	task := collect.Crawler{
		StartURL:       srv.URL,
		AllowedURL:     ".*",
		Depth:          1,
		EntitySelector: ".article--ssr",
		Extractor: func(*colly.HTMLElement, *goquery.Selection) error {
			dispatched = true
			return nil
		},
		Collector: collector,
	}

	err := task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestReal(t *testing.T) {

	dispatched := false

	task := collect.Crawler{
		StartURL:       "https://thailand-news.ru",
		AllowedURL:     "https://thailand-news\\.ru?.+",
		EntityURL:      "https://thailand-news\\.ru/news/((?:[^/]+/)*[^/]+)/.+",
		Depth:          3,
		EntitySelector: ".node-article--full",
		Extractor: func(c *colly.HTMLElement, q *goquery.Selection) error {
			dispatched = true
			println(c.Request.URL.String())
			return nil
		},
		Collector: nil,
	}
	err := task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestSave(t *testing.T) {

	srv := ServeFile(t, "run_test.html")
	defer srv.Close()

	dispatched := false

	dispatcher := func(payload *extract.Payload) error {
		dispatched = true
		return nil
	}

	extractor, err := store.NewExtractStore(testDbName, mongoCfg)
	require.NoError(t, err)

	task := collect.Crawler{
		StartURL:       srv.URL,
		AllowedURL:     ".*",
		Depth:          1,
		EntitySelector: ".article--ssr",
		Extractor:      extract.Pipe(spider.WindmillMeta, extract.Html, extract.Article, dispatcher, extractor.Save),
		Collector:      collector,
	}

	// expected ONE result, since we run Chromedp only if no results found
	// in this case first result is found in the HTML, so JS browse is not used
	err = task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestJSCollect(t *testing.T) {

	srv := ServeFile(t, "run_test.html")
	defer srv.Close()

	dispatched := false

	task := collect.Crawler{
		StartURL:       srv.URL,
		AllowedURL:     ".*",
		Depth:          1,
		UseBrowser:     true,
		EntitySelector: ".article--js",
		Extractor: func(*colly.HTMLElement, *goquery.Selection) error {
			dispatched = true
			return nil
		},
	}

	require.NoError(t, task.Start())
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

	config, err := mongodb.ConfigFromResource(validResource)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, config, "Expected non-nil config")
}

func TestServeFile(t *testing.T) {

	srv := ServeFile(t, "run_test.html")
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

type Args struct {
	// StartURL is the URL to start crawling, e.g. http://example.com
	StartURL string `json:"StartURL"`
	// AllowedURL is the regex to match the URLs, e.g. ".*"
	MatchURL string `json:"AllowedURL"`
	// Depth is the number of levels to follow the links
	Depth int `json:"Depth"`
	// Selector CSS to match the entities to extract, e.g. ".article--ssr"
	Selector string `json:"Selector"`
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

func TestParse(t *testing.T) {

	args := Args{}
	var input any = map[string]interface{}{
		"StartURL":   "http://example.com",
		"AllowedURL": ".*",
		"Depth":      1,
		"Selector":   ".article--ssr",
	}

	// Test the Parse function
	err := Parse(input, &args)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "http://example.com", args.StartURL)
	assert.Equal(t, ".*", args.MatchURL)
	assert.Equal(t, 1, args.Depth)
	assert.Equal(t, ".article--ssr", args.Selector)
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	return srv, nil
}

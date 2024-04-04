package spider_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/donq/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"spider"
	"spider/collect"
	"spider/extract"
	"spider/store"
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

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--ssr",
		Extract: func(*goquery.Selection, *url.URL) error {
			dispatched = true
			return nil
		},
		Storage: collector,
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

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--ssr",
		Extract:  extract.Pipe(spider.WindmillMeta, extract.Crawler, dispatcher, extractor.Save),
		Storage:  collector,
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

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--js",
		Extract: func(*goquery.Selection, *url.URL) error {
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
		w.Write(content)
	}))

	return srv, nil
}

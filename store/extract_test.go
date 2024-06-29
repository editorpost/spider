//go:build e2e
// +build e2e

package store

import (
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/payload"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	mongoCfg *mongodb.Config
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

	collector, err = NewCollectStore(testDbName, mongoCfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

//goland:noinspection GoUnusedFunction
func TestStart(t *testing.T) {

	storage, err := NewExtractStore("thailand-news_ru__full", mongoCfg.DSN)
	require.NoError(t, err)

	pool := proxy.NewPool("https://thailand-news.ru")
	require.NoError(t, pool.Start())

	task := collect.Crawler{
		StartURL:        "https://thailand-news.ru",
		AllowedURL:      "https://thailand-news.ru{any}",
		ExtractURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:           3,
		ExtractSelector: ".node-article--full",
		Extractor: payload.PipelineFn(extract.WindmillMeta, extract.Html, func(p *payload.Payload) error {
			println(p.URL.String())
			return nil
		}, storage.Save),
		Storage: nil,
		// ProxyFn:   pool.GetProxyURL,
		RoundTripper: pool.Transport(),
	}

	err = task.Run()
	require.NoError(t, err)
}

//goland:noinspection GoUnusedFunction
func TestStartNoDatabase(t *testing.T) {

	pool := proxy.NewPool("https://thailand-news.ru")
	require.NoError(t, pool.Start())

	task := collect.Crawler{
		StartURL:        "https://thailand-news.ru",
		AllowedURL:      "https://thailand-news.ru{any}",
		ExtractURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:           3,
		ExtractSelector: ".node-article--full",
		Extractor:       payload.PipelineFn(extract.WindmillMeta),
		Storage:         nil,
		RoundTripper:    pool.Transport(),
	}

	err = task.Run()
	require.NoError(t, err)
}

func TestDatabaseSave(t *testing.T) {

	srv := ServeFile(t, "local_data.html")
	defer srv.Close()

	dispatched := false

	dispatcher := func(payload *payload.Payload) error {
		dispatched = true
		return nil
	}

	extractor, err := NewExtractStore(testDbName, mongoCfg.DSN)
	require.NoError(t, err)

	task := collect.Crawler{
		StartURL:        srv.URL,
		AllowedURL:      ".*",
		Depth:           1,
		ExtractSelector: ".article--ssr",
		Extractor:       payload.PipelineFn(extract.WindmillMeta, extract.Html, dispatcher, extractor.Save),
		Storage:         collector,
	}

	// expected ONE result, since we run Chromedp only if no results found
	// in this case first result is found in the HTML, so JS browse is not used
	err = task.Run()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

//goland:noinspection GoUnusedFunction
func TestRealForbidden(t *testing.T) {

	storage, err := NewExtractStore("thailand-news_ru__full", mongoCfg.DSN)
	require.NoError(t, err)

	task := collect.Crawler{
		StartURL:        "https://thailand-news.ru/admin",
		AllowedURL:      "https://thailand-news.ru{any}",
		ExtractURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:           3,
		ExtractSelector: ".node-article--full",
		Extractor: payload.PipelineFn(extract.WindmillMeta, extract.Html, func(p *payload.Payload) error {
			println(p.URL.String())
			return nil
		}, storage.Save),
		Storage: nil,
	}

	err = task.Run()
	require.NoError(t, err)
}

var (
	collector *CollectStore
)

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

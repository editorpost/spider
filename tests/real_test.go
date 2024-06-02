//go:build e2e
// +build e2e

package tests

import (
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"github.com/stretchr/testify/require"
	"testing"
)

//goland:noinspection GoUnusedFunction
func TestStart(t *testing.T) {

	storage, err := store.NewExtractStore("thailand-news_ru__full", mongoCfg)
	require.NoError(t, err)

	pool := proxy.NewPool("https://thailand-news.ru")
	require.NoError(t, pool.Start())

	task := collect.Crawler{
		StartURL:       "https://thailand-news.ru",
		AllowedURL:     "https://thailand-news.ru{any}",
		EntityURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:          3,
		EntitySelector: ".node-article--full",
		Extractor: extract.Pipe(extract.WindmillMeta, extract.Html, func(p *extract.Payload) error {
			println(p.URL.String())
			return nil
		}, storage.Save),
		Storage: nil,
		// ProxyFn:   pool.GetProxyURL,
		RoundTripper: pool.Transport(),
	}

	err = task.Start()
	require.NoError(t, err)
}

//goland:noinspection GoUnusedFunction
func TestStartNoDatabase(t *testing.T) {

	pool := proxy.NewPool("https://thailand-news.ru")
	require.NoError(t, pool.Start())

	task := collect.Crawler{
		StartURL:       "https://thailand-news.ru",
		AllowedURL:     "https://thailand-news.ru{any}",
		EntityURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:          3,
		EntitySelector: ".node-article--full",
		Extractor:      extract.Pipe(extract.WindmillMeta),
		Storage:        nil,
		RoundTripper:   pool.Transport(),
	}

	err = task.Start()
	require.NoError(t, err)
}

func TestDatabaseSave(t *testing.T) {

	srv := ServeFile(t, "local_data.html")
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
		Extractor:      extract.Pipe(extract.WindmillMeta, extract.Html, dispatcher, extractor.Save),
		Storage:        collector,
	}

	// expected ONE result, since we run Chromedp only if no results found
	// in this case first result is found in the HTML, so JS browse is not used
	err = task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

//goland:noinspection GoUnusedFunction
func TestRealForbidden(t *testing.T) {

	storage, err := store.NewExtractStore("thailand-news_ru__full", mongoCfg)
	require.NoError(t, err)

	task := collect.Crawler{
		StartURL:       "https://thailand-news.ru/admin",
		AllowedURL:     "https://thailand-news.ru{any}",
		EntityURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:          3,
		EntitySelector: ".node-article--full",
		Extractor: extract.Pipe(extract.WindmillMeta, extract.Html, func(p *extract.Payload) error {
			println(p.URL.String())
			return nil
		}, storage.Save),
		Storage: nil,
	}

	err = task.Start()
	require.NoError(t, err)
}

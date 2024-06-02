package manage_test

import (
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/stretchr/testify/require"
	"testing"
)

//goland:noinspection GoUnusedFunction
func TestStartNoDatabase(t *testing.T) {

	pool := proxy.NewPool("https://thailand-news.ru")
	require.NoError(t, pool.Start())

	task := collect.Crawler{
		StartURL:       "https://thailand-news.ru/news/turizm/",
		AllowedURL:     "https://thailand-news.ru/news/{any}",
		EntityURL:      "https://thailand-news.ru/news/{dir}/{some}",
		Depth:          3,
		EntitySelector: ".node-article--full",
		Extractor:      extract.Pipe(extract.WindmillMeta),
		Storage:        nil,
		RoundTripper:   pool.Transport(),
	}

	err := task.Start()
	require.NoError(t, err)
}

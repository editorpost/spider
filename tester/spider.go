package tester

import (
	fk "github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewSpiderWith(t *testing.T, server *TestServer) *setup.Spider {

	args := NewArgs()
	args.StartURL = server.URL + "/index.html"
	args.AllowedURLs = []string{server.URL}
	args.ExtractURLs = []string{server.URL + "/{dir}/{some}-{num}.html"}
	args.ExtractLimit = 5
	args.Depth = 3

	s, err := setup.NewSpider(fk.UUID(), args, NewExtractConfig(), TestDeploy(t))
	require.NoError(t, err)
	require.NotNil(t, s)

	return s
}

func NewArgs() *config.Config {
	return &config.Config{
		StartURL:        "",
		ExtractSelector: "",
		ExtractLimit:    0,
		UseBrowser:      false,
		Depth:           0,
		UserAgent:       "",
		ProxyEnabled:    false,
		ProxySources:    nil,
	}
}

func NewExtractConfig() *extract.Config {
	return &extract.Config{
		Entities: []string{"article"},
		Fields: []*fields.Field{
			{
				Name:         "head_title",
				Cardinality:  1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head title",
			},
		},
		Media: &media.Config{
			Enabled: true,
		},
	}
}

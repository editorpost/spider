package tester

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewSpider(t *testing.T) *setup.Spider {
	s, err := setup.NewSpider(NewArgs(), NewExtractConfig())
	require.NoError(t, err)
	return s
}

func NewSpiderWith(t *testing.T, server *TestServer) *setup.Spider {

	args := NewArgs()
	args.StartURL = server.URL + "/index.html"
	args.AllowedURL = server.URL
	args.ExtractURL = server.URL + "/{dir}/{some}-{num}.html"
	args.ExtractLimit = 5
	args.Depth = 3

	s, err := setup.NewSpider(args, NewExtractConfig())
	require.NoError(t, err)
	require.NotNil(t, s)

	return s
}

func NewArgs() *config.Args {
	return &config.Args{
		ID:              gofakeit.UUID(),
		Name:            gofakeit.AppName(),
		StartURL:        "",
		AllowedURL:      "",
		ExtractURL:      "",
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
		ExtractEntities: []string{"article"},
		ExtractFields: []*fields.Field{
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

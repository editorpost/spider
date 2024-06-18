package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlags(t *testing.T) {

	argsJson, expectedArgs := NewArgs()
	treeJson, expectedTree := NewFields()

	flag.Set("args", argsJson)
	flag.Set("fields", treeJson)
	flag.Set("entities", "html,article")
	flag.Set("cmd", "start")

	cmd, args, entities, tree := Flags()

	assert.Equal(t, "start", cmd)
	assert.Equal(t, "html,article", entities)
	assert.Equal(t, expectedArgs, args)
	assert.Len(t, tree, len(expectedTree))
}

func NewArgs() (js string, args *config.Args) {
	args = &config.Args{
		SpiderID:        "ready-check",
		Name:            "Ready Check",
		StartURL:        "https://example.com",
		AllowedURL:      "https://example.com/{any}",
		ExtractURL:      "https://example.com/articles/{any}",
		ExtractSelector: "article",
		ExtractLimit:    1,
		UseBrowser:      true,
		Depth:           1,
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		ProxyEnabled:    true,
		ProxySources:    []string{},
	}

	// marshal args to JSON
	b, _ := json.Marshal(args)
	return string(b), args
}

func NewFields() (js string, f []*fields.Field) {

	f = []*fields.Field{
		{
			Name:         "title",
			Cardinality:  1,
			InputFormat:  "html",
			OutputFormat: []string{"text"},
			Selector:     ".product__title",
		},
		{
			Name:         "price",
			Cardinality:  1,
			InputFormat:  "html",
			OutputFormat: []string{"text"},
			Selector:     ".product__price--amount",
		},
	}

	// marshal fields to JSON
	b, _ := json.Marshal(f)
	return string(b), f
}
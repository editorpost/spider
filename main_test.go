package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFlags(t *testing.T) {

	jsonStr := `{
		"ID": "4c62e925-2e3e-4ff6-a92b-305b85d6281a",
		"StartURL": "http://example.com",
		"ExtractLimit": 10,
		"ExtractEntities": ["person", "organization"],
		"ExtractFields": [
			{
				"Name": "name",
				"Selector": "h1"
			}
		]
	}`

	s, err := setup.NewSpiderFromJSON([]byte(jsonStr))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal([]byte(jsonStr), s))

	require.NoError(t, flag.Set("spider", jsonStr))
	require.NoError(t, flag.Set("cmd", "start"))

	cmd, spider := Flags()

	assert.Equal(t, "start", cmd)
	assert.Equal(t, s.Args, spider.Args)
	assert.Equal(t, s.Config, spider.Config)
}

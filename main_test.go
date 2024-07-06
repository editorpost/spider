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
		"ID": "test",
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

	s := &setup.Spider{}
	require.NoError(t, json.Unmarshal([]byte(jsonStr), s))

	require.NoError(t, flag.Set("spider", jsonStr))
	require.NoError(t, flag.Set("cmd", "start"))

	cmd, spider := Flags()

	assert.Equal(t, "start", cmd)
	assert.Equal(t, s, spider)
}

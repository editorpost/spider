package config_test

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeps_Normalize(t *testing.T) {

	deps := config.Deps{}
	norm := deps.Normalize()
	assert.NotNil(t, norm.Extractor)
	assert.NotNil(t, norm.Monitor)
	assert.NotNil(t, norm.Storage)
}

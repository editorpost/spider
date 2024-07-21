package console_test

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrial(t *testing.T) {

	srv := tester.NewServer("../../tester/fixtures")
	defer srv.Close()

	s := tester.NewSpiderWith(t, srv)
	require.NotNil(t, s)

	data, err := console.Trial(s)
	require.NoError(t, err)
	assert.Len(t, data, 3)
}

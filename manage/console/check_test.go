package console_test

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheck(t *testing.T) {

	srv := tester.NewServer("../../tester/fixtures")
	defer srv.Close()

	s := tester.NewSpiderWith(t, srv)
	require.NotNil(t, s)

	_, err := console.Check(s)
	require.NoError(t, err)
}

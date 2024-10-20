package console_test

import (
	fk "github.com/brianvoe/gofakeit/v6"
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

	checkID := fk.UUID()
	id, err := console.Check(checkID, s)
	require.NoError(t, err)
	require.Equal(t, checkID, id)
}

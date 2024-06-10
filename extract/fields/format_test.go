package fields_test

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReduceSpaces(t *testing.T) {

	tc := []struct {
		name, in, out string
	}{
		{"empty", "", ""},
		{"spaces", "  ", ""},
		{"trim around", "  text  ", "text"},
		{"reduce spaces", "a          b", "a b"},
	}

	assert.NotNil(t, tc, "test cases are not defined")

	// use testify assert
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.out, fields.ReduceSpaces(c.in))
		})
	}
}

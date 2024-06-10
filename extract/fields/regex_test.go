package fields_test

import (
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestBetweenRegex(t *testing.T) {

	tc := []struct {
		name     string
		start    string
		end      string
		data     string
		expected []string
	}{
		{"empty", "", "", "", nil},
		{"simple", "Price:", "USD", "Price: 99.99 USD", []string{" 99.99 "}},
		{"not found", "Price:", "USD", "Price: 99.99", nil},
		{"multiple", "Price:", "USD", "Price: 99.99 USD Price: 89.99 USD", []string{" 99.99 ", " 89.99 "}},
		{"multiline", "Price:", "USD", "Price: 99.99 USD\nPrice: 89.99 USD", []string{" 99.99 ", " 89.99 "}},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {

			// empty start and end
			if c.start == "" && c.end == "" {
				assert.Nil(t, c.expected)
			}

			// compile and match
			between := regexp.MustCompile(regexp.QuoteMeta(c.start) + "(?s)(.*?)" + regexp.QuoteMeta(c.end))
			matches := fields.RegexExtract(between, c.data)
			assert.Equal(t, c.expected, matches)
		})
	}

}

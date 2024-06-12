package fields_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGroup(t *testing.T) {

	tc := []struct {
		name     string
		group    *fields.Group
		expected any
		hasErr   bool
		err      error
	}{
		{
			"single",
			&fields.Group{
				Name:     "product",
				Selector: ".product--full",
				Limit:    1,
				Required: true,
				Fields: []*fields.Field{
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
				},
			},
			map[string]any{
				"product": map[string]any{
					"title": "Main Product Title",
					"price": "99.99",
				},
			},
			false,
			nil,
		},
		{
			"multiple",
			&fields.Group{
				Name:     "products",
				Selector: ".product",
				Limit:    2,
				Required: true,
				Fields: []*fields.Field{
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
				},
			},
			map[string]any{
				"products": []any{
					map[string]any{
						"title": "Main Product Title",
						"price": "99.99",
					},
					map[string]any{
						"title": "Another Product Title",
						"price": "49.99",
					},
				},
			},
			false,
			nil,
		},
	}

	// use testify assert
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {

			// check error
			skipExpectedErr := func(actual error) bool {

				if actual == nil {
					// continue test case execution
					return false
				}

				// force error if not expected
				if !c.hasErr {
					assert.NoError(t, actual)
				}

				// check error instance
				if c.err != nil {
					assert.ErrorIs(t, c.err, actual)
				}

				// stops test case execution
				return true
			}

			fn, err := c.group.Extractor()
			if skip := skipExpectedErr(err); skip {
				return
			}

			// compare values
			read := strings.NewReader(GetTestFieldsHTML(t))
			dom, err := goquery.NewDocumentFromReader(read)
			require.NoError(t, err)

			values, err := fn(dom.Selection)
			if skip := skipExpectedErr(err); skip {
				return
			}

			assert.Equal(t, c.expected, values)
		})
	}
}

func TestGroupFromMap(t *testing.T) {

	m := map[string]any{
		"Name":     "product",
		"Selector": ".product--full",
		"Required": true,
		"Fields": []*fields.Field{
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
		},
	}

	e, err := fields.GroupFromMap(m)
	require.NoError(t, err)

	assert.Equal(t, "product", e.Name)
	assert.Equal(t, ".product--full", e.Selector)
	assert.True(t, e.Required)
	assert.Len(t, e.Fields, 2)

	title := e.Fields[0]
	assert.Equal(t, "title", title.Name)
	assert.Equal(t, 1, title.Cardinality)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := e.Fields[1]
	assert.Equal(t, "price", price.Name)
	assert.Equal(t, 1, price.Cardinality)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func TestGroupExtractorMap(t *testing.T) {

	e := &fields.Group{
		Name:     "product",
		Selector: ".product--full",
		Required: true,
		Fields: []*fields.Field{
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
		},
	}

	m := e.Map()

	assert.Equal(t, "product", m["Name"])
	assert.Equal(t, ".product--full", m["Selector"])
	assert.True(t, m["Required"].(bool))
	assert.Len(t, m["Fields"].([]*fields.Field), 2)

	title := m["Fields"].([]*fields.Field)[0]
	assert.Equal(t, "title", title.Name)
	assert.Equal(t, 1, title.Cardinality)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := m["Fields"].([]*fields.Field)[1]
	assert.Equal(t, "price", price.Name)
	assert.Equal(t, 1, price.Cardinality)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

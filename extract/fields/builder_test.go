package fields_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {

	tc := []struct {
		name     string
		builders []*fields.Field
		expected any
		hasErr   bool
		err      error
	}{
		{
			"single",
			[]*fields.Field{
				{
					Name:           "product",
					Selector:       ".product--full",
					Cardinality:    1,
					Required:       true,
					LimitSelection: true,
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
			[]*fields.Field{
				{
					Name:           "products",
					Selector:       ".product",
					Cardinality:    2,
					Required:       true,
					LimitSelection: true,
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
		{
			"multiple, skip missing required",
			[]*fields.Field{
				{
					Name:           "prices",
					Selector:       ".product__price",
					Cardinality:    0,
					Required:       true,
					LimitSelection: true,
					Fields: []*fields.Field{
						{
							Name:         "amount",
							Cardinality:  1,
							Required:     true,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--amount",
						},
						{
							Name:         "currency",
							Cardinality:  1,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--currency",
						},
					},
				},
			},
			map[string]any{
				"prices": []any{
					map[string]any{
						"amount":   "99.99",
						"currency": "USD",
					},
					map[string]any{
						"amount":   "49.99",
						"currency": "USD",
					},
					map[string]any{
						"amount":   "0.99",
						"currency": "USD",
					},
				},
			},
			false,
			nil,
		},
		{
			"multiple, skip missing required",
			[]*fields.Field{
				{
					Name:           "buy",
					Selector:       ".product__price",
					Cardinality:    1,
					Required:       true,
					LimitSelection: true,
					Fields: []*fields.Field{
						{
							Name:         "amount",
							Cardinality:  1,
							Required:     true,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--amount",
						},
						{
							Name:         "currency",
							Cardinality:  1,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--currency",
						},
					},
				},
				{
					Name:           "sell",
					Selector:       ".product__price",
					Cardinality:    2,
					Required:       true,
					LimitSelection: true,
					Fields: []*fields.Field{
						{
							Name:         "amount",
							Cardinality:  1,
							Required:     true,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--amount",
						},
						{
							Name:         "currency",
							Cardinality:  1,
							InputFormat:  "html",
							OutputFormat: []string{"text"},
							Selector:     ".product__price--currency",
						},
					},
				},
				{
					Name:         "title",
					Cardinality:  1,
					Required:     true,
					InputFormat:  "html",
					OutputFormat: []string{"text"},
					Selector:     "head title",
				},
			},
			map[string]any{
				"buy": map[string]any{
					"amount":   "99.99",
					"currency": "USD",
				},
				"sell": []any{
					map[string]any{
						"amount":   "99.99",
						"currency": "USD",
					},
					map[string]any{
						"amount":   "49.99",
						"currency": "USD",
					},
				},
				"title": "Product Page Example",
			},
			false,
			nil,
		},
		{
			"nil result",
			[]*fields.Field{
				{
					Name:           "prices",
					Selector:       ".product__price--amount",
					Cardinality:    1,
					Required:       true,
					LimitSelection: true,
					Fields: []*fields.Field{
						{
							Name:         "amount",
							Cardinality:  1,
							InputFormat:  "html",
							Required:     true,
							OutputFormat: []string{"text"},
							Selector:     ".product__price--amount",
						},
						{
							Name:         "currency",
							Cardinality:  1,
							InputFormat:  "html",
							Required:     true,
							OutputFormat: []string{"text"},
							Selector:     ".product__price--currency",
						},
					},
				},
			},
			nil,
			true,
			fields.ErrRequiredFieldMissing,
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

			fn, err := fields.Extract("", c.builders...)
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

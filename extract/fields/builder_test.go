package fields_test

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

var HTML *goquery.Document

func TestMain(m *testing.M) {
	HTML = GetTestDocument()
	m.Run()
}

type Case struct {
	name     string
	field    *fields.Field
	expected any
}

func TestConstructValidation(t *testing.T) {

	err := fields.Construct(&fields.Field{})
	assert.Error(t, err)

	var validatorErr validator.ValidationErrors
	assert.True(t, errors.As(err, &validatorErr))
	assert.Equal(t, len(validatorErr), 1)
}

func TestConstructRegex(t *testing.T) {

	field := &fields.Field{
		Name:       "product",
		FinalRegex: "[a-z",
		Required:   true,
	}

	err := fields.Construct(field)
	assert.Error(t, err)
}

func TestConstructRecursive(t *testing.T) {

	field := &fields.Field{
		Name: "product",
		Children: []*fields.Field{
			{
				Name:       "title",
				FinalRegex: "[a-z",
				Required:   true,
			},
		},
	}

	err := fields.Construct(field)
	assert.Error(t, err)
}

func TestExtract(t *testing.T) {

	tc := []Case{
		{
			"single",
			&fields.Field{

				Name:        "product",
				Selector:    ".product--full",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
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
		},
		{
			"fixed cardinality",
			&fields.Field{

				Name:        "products",
				Selector:    ".product",
				Cardinality: 2,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
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
		},
		{
			"skip required delta by field (no amount in 4rd product)",
			&fields.Field{

				Name:        "prices",
				Selector:    ".product__price",
				Cardinality: 0,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
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
		},
		{
			"multilevel, skip missing required",
			&fields.Field{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
					{
						Name:         "title",
						Cardinality:  1,
						Required:     true,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     "head title",
					},
					{
						Name:        "buy",
						Selector:    ".product__price",
						Cardinality: 1,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
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
						Name:        "sell",
						Selector:    ".product__price",
						Cardinality: 2,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
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
			},
			map[string]any{
				"offer": map[string]any{
					"title": "Product Page Example",
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
				},
			},
		},
		{
			"nested, deep invalid, not existing selector",
			&fields.Field{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
					{
						Name:         "title",
						Cardinality:  1,
						Required:     false,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     "head title",
					},
					{
						Name:        "sub-offer",
						Selector:    "",
						Cardinality: 1,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
							{
								Name:         "title",
								Cardinality:  1,
								Required:     false,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     "head title",
							},
							{
								Name:         "test-missing",
								Cardinality:  1,
								Required:     true,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     "not-a-selector-at-all",
							},
						},
					},
				},
			},
			nil,
		},
		{
			"nil result",
			&fields.Field{
				Name:        "prices",
				Selector:    ".product__price--amount",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
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
			nil,
		},
	}

	for _, c := range tc {
		t.Run(c.name, CaseHandler(c))
	}
}

func TestExtractNestedRequired(t *testing.T) {

	tc := []Case{
		{
			"skip extraction, nested, deep invalid, not existing selector",
			&fields.Field{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
					{
						Name:         "title",
						Cardinality:  1,
						Required:     false,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     "head title",
					},
					{
						Name:        "sub-offer",
						Selector:    "",
						Cardinality: 1,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
							{
								Name:         "title",
								Cardinality:  1,
								Required:     false,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     "head title",
							},
							{
								Name:         "test-missing",
								Cardinality:  1,
								Required:     true,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     "not-a-selector-at-all",
							},
						},
					},
				},
			},
			nil,
		},
		{
			"nested, deep valid, fixed selector",
			&fields.Field{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
					{
						Name:         "title",
						Cardinality:  1,
						Required:     false,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     "head title",
					},
					{
						Name:        "sub-offer",
						Selector:    "",
						Cardinality: 1,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
							{
								Name:         "title",
								Cardinality:  1,
								Required:     false,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     "head title",
							},
							{
								Name:         "test-valid",
								Cardinality:  1,
								Required:     true,
								InputFormat:  "html",
								OutputFormat: []string{"text"},
								Selector:     ".product--full .product__title",
							},
						},
					},
				},
			},
			map[string]any{
				"offer": map[string]any{
					"title": "Product Page Example",
					"sub-offer": map[string]any{
						"title":      "Product Page Example",
						"test-valid": "Main Product Title",
					},
				},
			},
		},

		{
			"partial valid, deep",
			&fields.Field{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Field{
					{
						Name:         "title",
						Cardinality:  1,
						Required:     true,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     "head title",
					},
					{
						Name:        "sell",
						Selector:    ".product__price",
						Cardinality: 2,
						Required:    true,
						Scoped:      true,
						Children: []*fields.Field{
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
			},
			map[string]any{
				"offer": map[string]any{
					"title": "Product Page Example",
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
				},
			},
		},
	}

	for _, c := range tc {
		t.Run(c.name, CaseHandler(c))
	}
}

func CaseHandler(c Case) func(t *testing.T) {
	return func(t *testing.T) {

		payload := map[string]any{}
		assert.NoError(t, fields.Construct(c.field))
		fields.Extract(payload, HTML.Selection, c.field)

		if c.expected == nil {
			assert.Empty(t, payload)
			return
		}

		assert.Equal(t, c.expected, payload)
	}
}

func BenchmarkExtract(b *testing.B) {

	field := &fields.Field{
		Name:        "product",
		Selector:    ".product--full",
		Cardinality: 1,
		Required:    true,
		Scoped:      true,
		Children: []*fields.Field{
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

	payload := map[string]any{}

	for i := 0; i < b.N; i++ {
		fields.Extract(payload, HTML.Selection, field)
	}
}

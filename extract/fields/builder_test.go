package fields_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"testing"
)

var HTML *goquery.Document

func TestMain(m *testing.M) {
	HTML = GetTestDocument()
	m.Run()
}

type Case struct {
	name      string
	extractor *fields.Extractor
	expected  any
	hasErr    bool
	err       error
}

func TestExtract(t *testing.T) {

	tc := []Case{
		{
			"single",
			&fields.Extractor{

				Name:        "product",
				Selector:    ".product--full",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
			&fields.Extractor{

				Name:        "products",
				Selector:    ".product",
				Cardinality: 2,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
		{
			"multiple, skip missing required",
			&fields.Extractor{

				Name:        "prices",
				Selector:    ".product__price",
				Cardinality: 0,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
			false,
			nil,
		},
		{
			"multiple, skip missing required",
			&fields.Extractor{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
						Children: []*fields.Extractor{
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
						Children: []*fields.Extractor{
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
			false,
			nil,
		},
		{
			"multiple, skip missing required",
			&fields.Extractor{
				Name:        "offer",
				Selector:    "",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
						Children: []*fields.Extractor{
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
			false,
			nil,
		},
		{
			"nil result",
			&fields.Extractor{
				Name:        "prices",
				Selector:    ".product__price--amount",
				Cardinality: 1,
				Required:    true,
				Scoped:      true,
				Children: []*fields.Extractor{
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
			true,
			fields.ErrRequiredFieldMissing,
		},
	}

	for _, c := range tc {
		t.Run(c.name, CaseHandler(c))
	}
}

func CaseHandler(c Case) func(t *testing.T) {

	return func(t *testing.T) {

		// every time recompile
		AssertErrorIs(t, c.hasErr, fields.Construct(c.extractor), c.err)

		payload := map[string]any{}
		err := fields.Extract(payload, HTML.Selection, c.extractor)
		AssertErrorIs(t, c.hasErr, err, c.err)

		if !c.hasErr {
			assert.Equal(t, c.expected, payload)
		}
	}
}

func AssertErrorIs(t *testing.T, expected bool, got error, want error) {

	t.Helper()

	if got == nil {
		// continue test case execution
		return
	}

	// force error if not expected
	if !expected {
		assert.NoError(t, got)
	}

	// check error instance
	if want != nil {
		assert.ErrorIs(t, got, want)
	}

	// stops test case execution
	return
}

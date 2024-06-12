package fields_test

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestErrorIs(t *testing.T) {
	var fail = errors.New("failed")
	err := fmt.Errorf("chain: %w", fail)

	assert.ErrorIs(t, err, fail)
}

func TestField(t *testing.T) {

	tc := []Case{
		{
			"empty",
			&fields.Extractor{},
			nil,
			true, // fExtractor name is required
			nil,
		},
		{
			"simple",
			&fields.Extractor{
				Name:         "title",
				Cardinality:  1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__title",
			},
			map[string]any{"title": "Main Product Title"},
			false,
			nil,
		},
		{
			"Between",
			&fields.Extractor{
				Name:         "image",
				Cardinality:  1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__price",
				BetweenStart: "Price:",
				BetweenEnd:   "USD",
			},
			map[string]any{"image": "99.99"},
			false,
			nil,
		},
		{
			"Between image from item prop",
			&fields.Extractor{
				Name:         "image",
				Cardinality:  1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			map[string]any{
				"image": "product-image.jpg",
			},
			false,
			nil,
		},
		{
			"Between multiple selections",
			&fields.Extractor{
				Name:         "images",
				Cardinality:  2,
				InputFormat:  "html",
				OutputFormat: []string{"html"},
				Selector:     "meta", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			map[string]any{
				"images": []any{
					"product-image.jpg",
				},
			},
			false,
			nil,
		},
		{
			"regex",
			&fields.Extractor{
				Name:         "category",
				Cardinality:  1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product--full",
				// multiline regex
				FinalRegex: "Category:(?s)(.*?)</p>",
			},
			map[string]any{"category": "Magic wands"},
			false,
			nil,
		},
		{
			"regex image from item prop",
			&fields.Extractor{
				Name:         "category",
				Cardinality:  1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head",
				FinalRegex:   "meta itemprop=\"image\" content=\"(.+?)\"",
			},
			map[string]any{"category": "product-image.jpg"},
			false,
			nil,
		},
		{
			"all prices",
			&fields.Extractor{
				Name:        "prices",
				Cardinality: 0,
				Selector:    ".product__price--amount",
			},
			map[string]any{"prices": []any{"99.99", "49.99", "0.99"}},
			false,
			nil,
		},
		{
			"all prices with limit",
			&fields.Extractor{
				Name:        "prices",
				Cardinality: 2,
				Selector:    ".product__price--amount",
			},
			map[string]any{"prices": []any{"99.99", "49.99"}},
			false,
			nil,
		},
		{
			"required fExtractor are empty",
			&fields.Extractor{
				Name:        "not-exists",
				Cardinality: 0,
				Selector:    ".product__not-exists-element",
				Required:    true,
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

func TestFieldFromMap(t *testing.T) {

	m := map[string]any{
		"Name":         "title",
		"Cardinality":  1,
		"InputFormat":  "html",
		"OutputFormat": []string{"text"},
		"Selector":     ".product__title",
	}

	e, err := fields.FieldFromMap(m)
	require.NoError(t, err)

	assert.Equal(t, "title", e.Name)
	assert.Equal(t, 1, e.Cardinality)
	assert.Equal(t, "html", e.InputFormat)
	assert.Equal(t, []string{"text"}, e.OutputFormat)
	assert.Equal(t, ".product__title", e.Selector)
}

func TestFieldToMap(t *testing.T) {

	e := &fields.Extractor{
		Name:         "title",
		Cardinality:  1,
		InputFormat:  "html",
		OutputFormat: []string{"text"},
		Selector:     ".product__title",
	}

	m := e.Map()

	assert.Equal(t, "title", m["Name"])
	assert.Equal(t, 1, m["Cardinality"])
	assert.Equal(t, "html", m["InputFormat"])
	assert.Equal(t, []string{"text"}, m["OutputFormat"])
	assert.Equal(t, ".product__title", m["Selector"])
}

func TestEntityTransformNewDocumentFromReaderError(t *testing.T) {
	extractor := &fields.Extractor{
		InputFormat:  "html",
		OutputFormat: []string{"text"},
	}
	input := "<div>Hello  world!</div>"
	output := fields.EntryTransform(extractor, input)
	assert.Equal(t, "Hello world!", output)
}

func TestGroup(t *testing.T) {

	tc := []Case{
		{
			"single",
			&fields.Extractor{
				Name:        "product",
				Selector:    ".product--full",
				Cardinality: 1,
				Required:    true,
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
	}

	for _, c := range tc {
		t.Run(c.name, CaseHandler(c))
	}
}

func TestGroupExtractorMap(t *testing.T) {

	e := &fields.Extractor{
		Name:     "product",
		Selector: ".product--full",
		Required: true,
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
	}

	m := e.Map()

	assert.Equal(t, "product", m["Name"])
	assert.Equal(t, ".product--full", m["Selector"])
	assert.True(t, m["Required"].(bool))
	assert.Len(t, m["Children"].([]*fields.Extractor), 2)

	title := m["Children"].([]*fields.Extractor)[0]
	assert.Equal(t, "title", title.Name)
	assert.Equal(t, 1, title.Cardinality)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := m["Children"].([]*fields.Extractor)[1]
	assert.Equal(t, "price", price.Name)
	assert.Equal(t, 1, price.Cardinality)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func GetTestFieldsHTML() string {

	b, err := os.ReadFile("field_test.html")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func GetTestDocument() *goquery.Document {

	code := GetTestFieldsHTML()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(code))
	if err != nil {
		panic(err)
	}

	return doc
}

//
//<!DOCTYPE html>
//<html lang="en">
//<head>
//<meta charset="UTF-8">
//<meta name="viewport" content="width=device-width, initial-scale=1.0">
//<title>Product Page Example</title>
//<-- Schema.org meta data -->
//<meta itemprop="name" content="Main Product Title">
//<meta itemprop="description" content="This is an example product description. It provides details about the product, its features, and benefits.">
//<meta itemprop="image" content="product-image.jpg">
//</head>
//<body>
//
//<!-- Single Product Section -->
//<div class="product product--full product-123">
//<h1 class="product__title">Main Product Title</h1>
//<div class="product__details">
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">99.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__category">
//Category: Magic wands
//</p>
//<p class="product__description">This is an example product description. It provides details about the product, its features, and benefits.</p>
//<ul class="product__features">
//<li class="product__feature">Feature 1</li>
//<li class="product__feature">Feature 2</li>
//<li class="product__feature">Feature 3</li>
//</ul>
//<div class="product__rating">
//<span class="product__rating-stars">★★★★☆</span>
//<span class="product__rating-count">(25 reviews)</span>
//</div>
//<button class="product__cart">Add to Cart</button>
//</div>
//</div>
//
//<!-- Multiple Products Section -->
//<div class="products">
//<div class="product product--related product-124">
//<h2 class="product__title">Another Product Title</h2>
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">49.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__description">Another product description providing essential details.</p>
//<button class="product__cart">Add to Cart</button>
//</div>
//
//<div class="product product--related product-125">
//<h2 class="product__title">Third Product Title</h2>
//<p class="product__price">
//<span class="product__price--label">Price:</span>
//<span class="product__price--amount">0.99</span>
//<span class="product__price--sale">USD</span>
//</p>
//<p class="product__description">A brief description of the third product.</p>
//<button class="product__cart">Add to Cart</button>
//</div>
//</div>
//
//<!-- Additional Information Section -->
//<div class="additional-info">
//<h2 class="additional-info__title">Additional Information</h2>
//<p class="additional-info__content">This section contains additional information related to the products, such as shipping details, return policy, and FAQs.</p>
//</div>
//</body>
//</html>

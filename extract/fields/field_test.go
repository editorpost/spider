package fields_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/fields"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildExtractor(t *testing.T) {

	tc := []struct {
		name      string
		extractor *fields.Field
		expected  any
		hasErr    bool
		err       error
	}{
		{
			"empty",
			&fields.Field{},
			nil,
			true, // field name is required
			nil,
		},
		{
			"simple",
			&fields.Field{
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__title",
			},
			"Main Product Title",
			false,
			nil,
		},
		{
			"between",
			&fields.Field{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__price",
				BetweenStart: "Price:",
				BetweenEnd:   "USD",
			},
			"99.99",
			false,
			nil,
		},
		{
			"between image from item prop",
			&fields.Field{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			"product-image.jpg",
			false,
			nil,
		},
		{
			"between multiple selections",
			&fields.Field{
				FieldName:    "muliple image",
				Limit:        2,
				InputFormat:  "html",
				OutputFormat: []string{"html"},
				Selector:     "meta", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			[]string{"product-image.jpg"},
			false,
			nil,
		},
		{
			"regex",
			&fields.Field{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product--full",
				// multiline regex
				FinalRegex: "Category:(?s)(.*?)</p>",
			},
			"Magic wands",
			false,
			nil,
		},
		{
			"regex image from item prop",
			&fields.Field{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head",
				FinalRegex:   "meta itemprop=\"image\" content=\"(.+?)\"",
			},
			"product-image.jpg",
			false,
			nil,
		},
		{
			"all prices",
			&fields.Field{
				FieldName: "prices",
				Limit:     0,
				Selector:  ".product__price--amount",
			},
			[]string{"99.99", "49.99", "0.99"},
			false,
			nil,
		},
		{
			"all prices with limit",
			&fields.Field{
				FieldName: "prices",
				Limit:     2,
				Selector:  ".product__price--amount",
			},
			[]string{"99.99", "49.99"},
			false,
			nil,
		},
		{
			"required field are empty",
			&fields.Field{
				FieldName: "not-exists",
				Limit:     0,
				Selector:  ".product__not-exists-element",
				Required:  true,
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

			fn, err := c.extractor.Extractor()
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

func TestBuildGroup(t *testing.T) {

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
				Fields: map[string]*fields.Field{
					"title": {
						FieldName:    "title",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__title",
					},
					"price": {
						FieldName:    "price",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__price--amount",
					},
				},
			},
			map[string]any{
				"title": "Main Product Title",
				"price": "99.99",
			},
			false,
			nil,
		},
		{
			"multiple",
			&fields.Group{
				Name:     "products",
				Selector: ".product",
				Limit:    1,
				Required: true,
				Fields: map[string]*fields.Field{
					"title": {
						FieldName:    "title",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__title",
					},
					"price": {
						FieldName:    "price",
						Limit:        1,
						InputFormat:  "html",
						OutputFormat: []string{"text"},
						Selector:     ".product__price--amount",
					},
				},
			},
			map[string]any{
				"title": "Main Product Title",
				"price": "99.99",
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

func TestExtractorFromMap(t *testing.T) {

	m := map[string]any{
		"FieldName":    "title",
		"Limit":        1,
		"InputFormat":  "html",
		"OutputFormat": []string{"text"},
		"Selector":     ".product__title",
	}

	e, err := fields.FieldFromMap(m)
	require.NoError(t, err)

	assert.Equal(t, "title", e.FieldName)
	assert.Equal(t, 1, e.Limit)
	assert.Equal(t, "html", e.InputFormat)
	assert.Equal(t, []string{"text"}, e.OutputFormat)
	assert.Equal(t, ".product__title", e.Selector)
}

func TestGroupFromMap(t *testing.T) {

	m := map[string]any{
		"Name":     "product",
		"Selector": ".product--full",
		"Required": true,
		"Fields": map[string]*fields.Field{
			"title": {
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__title",
			},
			"price": {
				FieldName:    "price",
				Limit:        1,
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

	title := e.Fields["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := e.Fields["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func TestExtractorMap(t *testing.T) {

	e := &fields.Field{
		FieldName:    "title",
		Limit:        1,
		InputFormat:  "html",
		OutputFormat: []string{"text"},
		Selector:     ".product__title",
	}

	m := e.Map()

	assert.Equal(t, "title", m["FieldName"])
	assert.Equal(t, 1, m["Limit"])
	assert.Equal(t, "html", m["InputFormat"])
	assert.Equal(t, []string{"text"}, m["OutputFormat"])
	assert.Equal(t, ".product__title", m["Selector"])
}

func TestGroupExtractorMap(t *testing.T) {

	e := &fields.Group{
		Name:     "product",
		Selector: ".product--full",
		Required: true,
		Fields: map[string]*fields.Field{
			"title": {
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product__title",
			},
			"price": {
				FieldName:    "price",
				Limit:        1,
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
	assert.Len(t, m["Fields"].(map[string]*fields.Field), 2)

	title := m["Fields"].(map[string]*fields.Field)["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := m["Fields"].(map[string]*fields.Field)["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func TestEntityTransformNewDocumentFromReaderError(t *testing.T) {
	extractor := &fields.Field{
		InputFormat:  "html",
		OutputFormat: []string{"text"},
	}
	input := "<div>Hello  world!</div>"
	output := fields.EntryTransform(extractor, input)
	assert.Equal(t, "Hello world!", output)
}

func GetTestFieldsHTML(t *testing.T) string {

	t.Helper()

	// open file `article_test.html` return as string
	f, err := os.Open("fields_test.html")
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
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

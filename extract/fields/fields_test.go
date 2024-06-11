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
		extractor *fields.Extractor
		expected  []any
		hasErr    bool
		err       error
	}{
		{
			"empty",
			&fields.Extractor{},
			nil,
			true, // field name is required
			nil,
		},
		{
			"simple",
			&fields.Extractor{
				FieldName:    "title",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__title",
			},
			[]any{"Main Product Title"},
			false,
			nil,
		},
		{
			"between",
			&fields.Extractor{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "text",
				OutputFormat: []string{"text"},
				Selector:     ".product--full .product__price",
				BetweenStart: "Price:",
				BetweenEnd:   "USD",
			},
			[]any{"99.99"},
			false,
			nil,
		},
		{
			"between image from item prop",
			&fields.Extractor{
				FieldName:    "image",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"between multiple selections",
			&fields.Extractor{
				FieldName:    "muliple image",
				Limit:        10,
				InputFormat:  "html",
				OutputFormat: []string{"html"},
				Selector:     "meta", // multiple selection
				BetweenStart: "itemprop=\"image\" content=\"",
				BetweenEnd:   "\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"regex",
			&fields.Extractor{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     ".product--full",
				// multiline regex
				FinalRegex: "Category:(?s)(.*?)</p>",
			},
			[]any{"Magic wands"},
			false,
			nil,
		},
		{
			"regex image from item prop",
			&fields.Extractor{
				FieldName:    "category",
				Limit:        1,
				InputFormat:  "html",
				OutputFormat: []string{"text"},
				Selector:     "head",
				FinalRegex:   "meta itemprop=\"image\" content=\"(.+?)\"",
			},
			[]any{"product-image.jpg"},
			false,
			nil,
		},
		{
			"all prices",
			&fields.Extractor{
				FieldName: "prices",
				Limit:     0,
				Selector:  ".product__price--amount",
			},
			[]any{"99.99", "49.99", "0.99"},
			false,
			nil,
		},
		{
			"all prices with limit",
			&fields.Extractor{
				FieldName: "prices",
				Limit:     2,
				Selector:  ".product__price--amount",
			},
			[]any{"99.99", "49.99"},
			false,
			nil,
		},
		{
			"required field are empty",
			&fields.Extractor{
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

			fn, err := fields.BuildExtractor(c.extractor)
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
		group    *fields.GroupExtractor
		expected any
		hasErr   bool
		err      error
	}{
		{
			"single",
			&fields.GroupExtractor{
				Name:     "product",
				Selector: ".product--full",
				Required: true,
				Extractors: map[string]*fields.Extractor{
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
			[]any{
				map[string]any{
					"title": []any{"Main Product Title"},
					"price": []any{"99.99"},
				},
			},
			false,
			nil,
		},
		//{
		//	"multiple",
		//	&fields.GroupExtractor{
		//		Name:     "product",
		//		Selector: ".product",
		//		Required: true,
		//		Extractors: map[string]*fields.Extractor{
		//			"title": {
		//				FieldName:    "title",
		//				Limit:        1,
		//				InputFormat:  "html",
		//				OutputFormat: []string{"text"},
		//				Selector:     ".product__title",
		//			},
		//			"price": {
		//				FieldName:    "price",
		//				Limit:        1,
		//				InputFormat:  "html",
		//				OutputFormat: []string{"text"},
		//				Selector:     ".product__price--amount",
		//			},
		//		},
		//	},
		//	[]map[string]any{
		//		{
		//			"title": []any{"Main Product Title"},
		//			"price": []any{"99.99"},
		//		},
		//		{
		//			"title": []any{"Another Product Title"},
		//			"price": []any{"49.99"},
		//		},
		//		{
		//			"title": []any{"Third Product Title"},
		//			"price": []any{"0.99"},
		//		},
		//	},
		//	false,
		//	nil,
		//},
		//{
		//	"required field are empty",
		//	&fields.GroupExtractor{
		//		Name:     "product",
		//		Selector: ".product--not-exists",
		//		Required: true,
		//		Extractors: map[string]*fields.Extractor{
		//			"title": {
		//				FieldName:    "title",
		//				Limit:        1,
		//				InputFormat:  "html",
		//				OutputFormat: []string{"text"},
		//				Selector:     ".product__title",
		//			},
		//			"price": {
		//				FieldName:    "price",
		//				Limit:        1,
		//				InputFormat:  "html",
		//				OutputFormat: []string{"text"},
		//				Selector:     ".product__price--amount",
		//			},
		//		},
		//	},
		//	nil,
		//	true,
		//	fields.ErrRequiredFieldMissing,
		//},
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

			fn, err := fields.BuildGroup(c.group)
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

	e, err := fields.ExtractorFromMap(m)
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
		"Extractors": map[string]*fields.Extractor{
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
	assert.Len(t, e.Extractors, 2)

	title := e.Extractors["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := e.Extractors["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
}

func TestExtractorMap(t *testing.T) {

	e := &fields.Extractor{
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

	e := &fields.GroupExtractor{
		Name:     "product",
		Selector: ".product--full",
		Required: true,
		Extractors: map[string]*fields.Extractor{
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
	assert.Len(t, m["Extractors"].(map[string]*fields.Extractor), 2)

	title := m["Extractors"].(map[string]*fields.Extractor)["title"]
	assert.Equal(t, "title", title.FieldName)
	assert.Equal(t, 1, title.Limit)
	assert.Equal(t, "html", title.InputFormat)
	assert.Equal(t, []string{"text"}, title.OutputFormat)
	assert.Equal(t, ".product__title", title.Selector)

	price := m["Extractors"].(map[string]*fields.Extractor)["price"]
	assert.Equal(t, "price", price.FieldName)
	assert.Equal(t, 1, price.Limit)
	assert.Equal(t, "html", price.InputFormat)
	assert.Equal(t, []string{"text"}, price.OutputFormat)
	assert.Equal(t, ".product__price--amount", price.Selector)
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

func TestEntityTransformNewDocumentFromReaderError(t *testing.T) {
	extractor := &fields.Extractor{
		InputFormat:  "html",
		OutputFormat: []string{"text"},
	}
	input := "<div>Hello  world!</div>"
	output := fields.EntryTransform(extractor, input)
	assert.Equal(t, "Hello world!", output)
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

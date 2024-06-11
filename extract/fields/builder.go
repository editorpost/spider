package fields

import "github.com/PuerkitoBio/goquery"

type (
	ExtractFn func(*goquery.Selection) (any, error)

	Builder interface {
		Extractor() (ExtractFn, error)
	}
)

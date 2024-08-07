package extract

import (
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/extract/pipe"
	"log/slog"
)

//goland:noinspection GoNameStartsWithPackageName

func Extractors(ff []*fields.Field, entities ...string) ([]pipe.Extractor, error) {

	// entity extractors
	extractors := ExtractorsByName(entities...)

	// field extractors
	if len(ff) > 0 {
		extractFields, err := Fields(ff...)
		if err != nil {
			slog.Error("build extractors from field tree", slog.String("err", err.Error()))
			return nil, err
		}
		extractors = append(extractors, extractFields)
	}

	return extractors, nil
}

// ExtractorsByName creates slice of extractors by name.
// The string is a string like "html,article", e.g.: extract.Html, extract.Article
func ExtractorsByName(names ...string) []pipe.Extractor {

	if len(names) == 0 {
		return []pipe.Extractor{}
	}

	extractors := make([]pipe.Extractor, 0)
	for _, key := range names {
		switch key {
		case "html":
			extractors = append(extractors, Html)
		case "article":
			extractors = append(extractors, article.Article)
		}
	}

	return extractors
}

// ExtractorsByJsonString creates slice of extractors by name.
//
//goland:noinspection GoUnusedExportedFunction
func ExtractorsByJsonString(js string) []pipe.Extractor {
	if js == "" {
		return []pipe.Extractor{}
	}
	return []pipe.Extractor{}
}

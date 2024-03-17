package extract

import (
	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"log/slog"
	"net/url"
)

func Process(node *html.Node, u *url.URL) error {

	article, err := readability.FromDocument(node, u)
	if err != nil {
		slog.Error("extract failed", err)
		return err
	}
	slog.Debug("extract success", slog.String("title", article.Title))

	return nil
}

package extract

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"strings"
)

func Html(p *Payload) (err error) {
	p.Data[HtmlField], err = p.Doc.DOM.Html()
	return err
}

// LoadHTMLFromFile and returns it as string.
func LoadHTMLFromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CreatePayload from HTML string and URL.
func CreatePayload(htmlStr string, urlStr string) (*Payload, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Payload{
		Selection: doc.Selection,
		URL:       u,
		Data:      make(map[string]any),
	}, nil
}

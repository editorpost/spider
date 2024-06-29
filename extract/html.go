package extract

import "github.com/editorpost/spider/extract/payload"

func Html(p *payload.Payload) (err error) {
	p.Data[payload.HtmlField], err = p.Doc.DOM.Html()
	return err
}

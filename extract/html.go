package extract

import "github.com/editorpost/spider/extract/pipe"

func Html(p *pipe.Payload) (err error) {
	p.Data[pipe.HtmlField], err = p.Doc.DOM.Html()
	return err
}

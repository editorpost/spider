package extract

func Html(p *Payload) (err error) {
	p.Data[HtmlField], err = p.Doc.DOM.Html()
	return err
}

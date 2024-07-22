package media

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/pipe"
	"log/slog"
	"net/url"
	"strings"
	"sync"
)

// Claim claim to get src media and save it to dst.
type Claim struct {
	// Src is the Endpoint of the media to download.
	Src string `json:"Src"`
	// Dst is the path to save the downloaded media.
	Dst string `json:"Dst"`
	// Requested is true if the media is requested to download.
	Requested bool `json:"Requested"`
	// Done is true if the media is downloaded to destination.
	Done bool `json:"Done"`
}

type Claims struct {
	// publicURL is a prefix of the replaced media url.
	publicURL string
	// claims keyed with source url
	claims map[string]*Claim
	mute   *sync.Mutex
}

// NewClaims creates a new Claim for each image and replace src path in document and selection.
// Replacement path from media.Filename. Replaces src url in selection.
func NewClaims(publicURL string) *Claims {

	claims := &Claims{
		publicURL: publicURL,
		claims:    make(map[string]*Claim),
		mute:      &sync.Mutex{},
	}

	return claims
}

// ExtractAndReplace Claim for each img tag and replace src path in selection.
func (list *Claims) ExtractAndReplace(payload *pipe.Payload) *Claims {

	selection := payload.Doc.DOM

	// selection is nil
	if selection == nil {
		slog.Error("extracting media claims: payload DOM selection is nil")
		return list
	}

	// og:image
	selection.Find("meta").Each(func(i int, el *goquery.Selection) {

		key, _ := el.Attr("property")
		if key != "og:image" {
			return
		}

		src, _ := el.Attr("content")
		src = payload.Doc.Request.AbsoluteURL(src)

		if dst, err := list.Add(src); err != nil {
			el.SetAttr("content", dst)
		}
	})

	// img tag
	selection.Find("img").Each(func(_ int, el *goquery.Selection) {

		// has src
		src, exists := el.Attr("src")
		if !exists {
			return
		}

		// add claim
		src = payload.Doc.Request.AbsoluteURL(src)
		// replace src with claimed destination
		if dst, err := list.Add(src); err != nil {
			el.SetAttr("src", dst)
		}
	})

	return list
}

func (list *Claims) Add(src string) (string, error) {

	// already replaced
	if strings.HasPrefix(src, list.publicURL) {
		return src, nil
	}

	// filename as src url hash
	filename, err := Filename(src)
	if err != nil {
		slog.Error("failed to hash filename", slog.String("src", src), slog.String("err", err.Error()))
		return "", err
	}

	// full url
	dst, err := url.JoinPath(list.publicURL, filename)
	if err != nil {
		slog.Error("failed to join url", slog.String("dst", list.publicURL), slog.String("filename", filename), slog.String("err", err.Error()))
		return "", err
	}

	c := &Claim{
		Src: src,
		Dst: dst,
	}

	list.mute.Lock()
	defer list.mute.Unlock()
	list.claims[src] = c

	return c.Dst, nil
}

// Request Claim for uploading by Dst url.
func (list *Claims) Request(uri string) *Claims {

	list.mute.Lock()
	defer list.mute.Unlock()

	if claim := list.byDestination(uri); claim != nil {
		claim.Requested = true
		return list
	}

	if claim := list.bySource(uri); claim != nil {
		claim.Requested = true
	}

	return list
}

func (list *Claims) byDestination(u string) *Claim {
	for _, claim := range list.claims {
		if claim.Dst == u {
			return claim
		}
	}
	return nil
}

func (list *Claims) bySource(u string) *Claim {
	for _, claim := range list.claims {
		if claim.Src == u {
			return claim
		}
	}
	return nil
}

// Done Claim marks Claim as uploaded.
func (list *Claims) Done(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			list.mute.Lock()
			claim.Done = true
			list.mute.Unlock()
			return list
		}
	}

	return list
}

// Requested returns a list of requested claims.
func (list *Claims) Requested() []*Claim {

	requested := make([]*Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Requested {
			requested = append(requested, claim)
		}
	}
	return requested
}

// Uploaded returns a list of uploaded claims.
func (list *Claims) Uploaded() []*Claim {

	uploaded := make([]*Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Done {
			uploaded = append(uploaded, claim)
		}
	}
	return uploaded
}

func (list *Claims) All() []*Claim {

	all := make([]*Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		all = append(all, claim)
	}
	return all
}

func (list *Claims) Len() int {
	return len(list.claims)
}

func AbsoluteUrl(base *url.URL, href string) string {

	// parse the href
	rel, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// already absolute
	if rel.Scheme != "" {
		return rel.String()
	}

	// resolve the base with the relative href
	abs := base.ResolveReference(rel)

	return abs.String()
}

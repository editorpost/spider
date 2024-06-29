package media

import (
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"path"
	"strings"
)

// Claim claim to get src media and save it to dst.
type Claim struct {
	// Src is the URL of the media to download.
	Src string `json:"Src"`
	// Dst is the path to save the downloaded media.
	Dst string `json:"Dst"`
	// Requested is true if the media is requested to download.
	Requested bool `json:"Requested"`
	// Done is true if the media is downloaded to destination.
	Done bool `json:"Done"`
}

type Claims struct {
	// dstURL is a prefix of public path of the replaced media url.
	dstURL string
	// claims keyed with source url
	claims map[string]Claim
}

// NewClaims creates a new Claim for each image and replace src path in document and selection.
// Replacement path from media.Filename. Replaces src url in selection.
func NewClaims(uri string) *Claims {

	claims := &Claims{
		dstURL: uri,
		claims: make(map[string]Claim),
	}

	return claims
}

// ExtractAndReplace Claim for each img tag and replace src path in selection.
func (list *Claims) ExtractAndReplace(uri string, selection *goquery.Selection) {

	selection.Find("img").Each(func(i int, el *goquery.Selection) {

		// has src
		src, exists := el.Attr("src")
		if !exists {
			return
		}

		// already claimed
		if _, exists = list.claims[src]; exists {
			return
		}

		// already replaced
		if strings.HasPrefix(src, list.dstURL) {
			return
		}

		// filename as src url hash
		dst, err := Filename(src)
		if err != nil {
			slog.Error("failed to hash filename", slog.String("src", src), slog.String("err", err.Error()))
			return
		}

		// full url
		dst = path.Join(uri, dst)

		// replace url in selection
		el.SetAttr("src", dst)

		// add claim
		list.Add(Claim{
			Src: src,
			Dst: dst,
		})
	})
}

func (list *Claims) Add(c Claim) *Claims {
	list.claims[c.Src] = c
	return list
}

// Request Claim for uploading by Dst url.
func (list *Claims) Request(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			claim.Requested = true
		}
	}

	return list
}

// Done Claim marks Claim as uploaded.
func (list *Claims) Done(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			claim.Done = true
		}
	}

	return list
}

// Requested returns a list of requested claims.
func (list *Claims) Requested() []Claim {

	requested := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Requested {
			requested = append(requested, claim)
		}
	}
	return requested
}

// Uploaded returns a list of uploaded claims.
func (list *Claims) Uploaded() []Claim {

	uploaded := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Done {
			uploaded = append(uploaded, claim)
		}
	}
	return uploaded
}

func (list *Claims) All() []Claim {

	all := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		all = append(all, claim)
	}
	return all
}

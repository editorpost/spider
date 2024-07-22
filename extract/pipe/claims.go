package pipe

import (
	"net/url"
	"strings"
	"sync"
)

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

func (list *Claims) Add(src string) (string, error) {

	// already replaced
	if strings.HasPrefix(src, list.publicURL) {
		return src, nil
	}

	// already claimed
	if claim := list.bySource(src); claim != nil {
		return claim.Dst, nil
	}

	// filename as src url hash
	filename, err := Filename(src)
	if err != nil {
		return "", err
	}

	// full url
	dst, err := url.JoinPath(list.publicURL, filename)
	if err != nil {
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

func (list *Claims) Empty() bool {
	return list.Len() == 0
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

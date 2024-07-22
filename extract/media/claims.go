package media

import (
	"errors"
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
	// DstPath relative path to media downloaded
	DstPath string `json:"DstPath"`
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
	mute   *sync.RWMutex
}

// NewClaims creates a new Claim for each image and replace src path in document and selection.
// Replacement path from media.Filename. Replaces src url in selection.
func NewClaims(publicURL string) *Claims {

	claims := &Claims{
		publicURL: publicURL,
		claims:    make(map[string]*Claim),
		mute:      &sync.RWMutex{},
	}

	return claims
}

func (list *Claims) Add(payloadID string, src string) (Claim, error) {

	// already replaced
	if strings.HasPrefix(src, list.publicURL) {
		// look by destination url
		if claim := list.byDestinationUrl(src); claim != nil {
			return *claim, nil
		}
		return Claim{}, errors.New("parsed url replaced, but not found in claims")
	}

	// already claimed
	if claim := list.bySourceUrl(src); claim != nil {
		return *claim, nil
	}

	// filename as src url hash
	filename, err := Filename(src)
	if err != nil {
		return Claim{}, err
	}

	// destination path
	dstPath, err := url.JoinPath(payloadID, "media", filename)
	if err != nil {
		return Claim{}, err
	}

	// destination url
	dstUrl, err := url.JoinPath(list.publicURL, dstPath)
	if err != nil {
		return Claim{}, err
	}

	c := &Claim{
		Src:     src,
		Dst:     dstUrl,
		DstPath: dstPath,
	}

	list.add(c)

	return *c, nil
}

// Done Claim marks Claim as uploaded.
func (list *Claims) Done(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			return list.done(claim)
		}
	}

	return list
}

// Uploaded returns a list of uploaded claims.
func (list *Claims) Uploaded() []*Claim {

	uploaded := make([]*Claim, 0, len(list.claims))
	for _, claim := range list.All() {
		if claim.Done {
			uploaded = append(uploaded, claim)
		}
	}
	return uploaded
}

func (list *Claims) bySourceUrl(u string) *Claim {

	for _, claim := range list.All() {
		if claim.Src == u {
			return claim
		}
	}
	return nil
}

func (list *Claims) byDestinationUrl(u string) *Claim {

	for _, claim := range list.All() {
		if claim.Dst == u {
			return claim
		}
	}
	return nil
}

func (list *Claims) Empty() bool {
	return list.Len() == 0
}

func (list *Claims) All() []*Claim {

	list.mute.RLock()
	defer list.mute.RUnlock()

	all := make([]*Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		all = append(all, claim)
	}
	return all
}

func (list *Claims) Len() int {

	list.mute.RLock()
	defer list.mute.RUnlock()

	return len(list.claims)
}

func (list *Claims) add(c *Claim) {
	list.mute.Lock()
	defer list.mute.Unlock()
	list.claims[c.Src] = c
}

func (list *Claims) done(c *Claim) *Claims {
	list.mute.Lock()
	defer list.mute.Unlock()
	c.Done = true
	return list
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

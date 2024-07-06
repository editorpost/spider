package proxy

// Source is the structure for the source list
// JSON representation of the source list:
//
//	{
//		"Kind": "uri",
//		"Endpoint": "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&anonymity=Elite&timeout=20000",
//		"Schema": "http"
//	}
type Source struct {
	// Kind of structure for the source list
	// Only `uri` expected right now - extend kind when needed
	Kind string
	// Endpoint of the source list
	Endpoint string
	// Schema for proxy IP if not provided, default is http
	// Expected one of http|https|socks4|socks5
	Schema string
}

// LoadPublicLists loads the valid list from public sources
func LoadPublicLists() ([]string, error) {
	return LoadStringLists([]string{
		// uri
		"https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&anonymity=Elite&timeout=20000",
		// host
		"https://sunny9577.github.io/proxy-scraper/proxies.txt",
		// host
		"https://www.proxy-list.download/api/v1/get?type=http",
	})
}

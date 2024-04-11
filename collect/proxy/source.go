package proxy

// LoadPublicList loads the valid list from public sources
func LoadPublicLists() ([]string, error) {
	return LoadStringLists([]string{
		"https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&anonymity=Elite&timeout=20000",
		"https://sunny9577.github.io/proxy-scraper/proxies.txt",
		"https://www.proxy-list.download/api/v1/get?type=http",
	})
}

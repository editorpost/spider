package proxy

// LoadPublicList loads the valid list from public sources
func LoadPublicLists() ([]string, error) {
	return LoadStringLists([]string{
		"https://api.proxyscrape.com/v1/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all",
		"https://sunny9577.github.io/proxy-scraper/proxies.txt",
		"https://www.proxy-list.download/api/v1/get?type=http",
	})
}

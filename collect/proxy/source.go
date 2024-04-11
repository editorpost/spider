package proxy

// LoadProxyScrapeList loads the valid list from proxyscrape.com
func LoadProxyScrapeList() ([]string, error) {
	uri := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&timeout=20000"
	return LoadStringList(uri)
}

// LoadSunnyProxyList loads the valid list from sunny9577.github.io
func LoadSunnyProxyList() ([]string, error) {
	uri := "https://sunny9577.github.io/proxy-scraper/proxies.txt"
	return LoadStringList(uri)
}

package collect_test

import (
	"github.com/editorpost/spider/collect"
	"testing"
)

func TestLoadProxyList(t *testing.T) {
	url := "https://sunny9577.github.io/proxy-scraper/proxies.json" // https://www.proxy-list.download/api/v1/get?type=http"
	_ = collect.LoadProxyList(url)
}

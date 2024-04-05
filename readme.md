
# Usage as Windmill Script

Expected arguments:
```json 
{
  "StartURL": "https://thailand-news.ru/news/turizm/khochu-na-pkhuket",
  "AllowedURL": "https://thailand-news.ru/news/.+",
  "EntityURL": "https://thailand-news.ru/news/tourizm/.+",
  "UseBrowser": false,
  "Depth": 1,
  "Selector": ".node-article--full"
}
```

Start task example:
```go 
package inner

import (
	"github.com/editorpost/spider"
)

func main(crawler interface{}) (interface{}, error) {
	//require github.com/editorpost/spider v0.0.0-20240405160747-44ad051bec9e
	return 0, spider.StartWith(crawler)
}

```

Init script for workers in Windmill:
```bash
# the last line of the stdout is the return value
# unless you write json to './result.json' or a string to './result.out'
echo "Hello $msg"
go clean -cache
go clean -modcache
rm -rf /tmp/windmill/cache/gobin/*
rm -rf /tmp/windmill/cache/go/pkg/mod/github.com/editorpost/
```
use it as comment in script to pin version
//require github.com/editorpost/spider 
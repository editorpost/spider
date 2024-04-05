
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
	//require github.com/editorpost/spider v0.0.0-20240405201109-79b03452a487
	return 0, spider.StartWith(crawler)
}

```

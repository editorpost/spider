
# Usage as Windmill Script

Expected arguments:
```json 
{
  "StartURL": "https://thailand-news.ru/news/turizm/khochu-na-pkhuket",
  "MatchURL": "https://thailand-news.ru/news/turizm/.+",
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
	return 0, spider.StartWith(crawler)
}

```
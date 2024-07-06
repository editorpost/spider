
# Disclaimer. Don't use it. Under development.


# Run as binary
```bash
go run main.go -cmd="start" -spider="{}"
```

# Usage as Windmill Script
Ensure you have the Windmill Mongodb resource `f/spider/resource/deploy` available in your Windmill environment.

Expected arguments:
```json 
{
  "StartURL": "https://news.com/news/business/new-offshore-zones-promoted",
  "AllowedURL": "https://news.com{any}",
  "EntityURL": "https://news.com/news/{dir}/{some}",
  "EntitySelector": ".node-article--full",
  "UseBrowser": false,
  "Depth": 1
}
```

### Start task example:
```go 
package inner

import (
	"github.com/editorpost/spider"
)

func main(crawler interface{}) (interface{}, error) {
	//require github.com/editorpost/spider v0.0.7
	return 0, spider.StartWith(crawler)
}
```

Note! You can use the specific version of the library providing comment with `//require repo/pkg v0.0.7`. The version is specified for example only.
Use `git rev-parse HEAD` to get the current version of the library.

### Initialization with user parameters:
```go
package inner

import (
	"github.com/editorpost/spider"
)

func main(
	name string,
	startURL string,
	allowedURL string,
	entityURL string,
	entitySelector string,
	depth int,
	useBrowser bool,
) (interface{}, error) {

	//require github.com/editorpost/spider v0.0.7
	args := &spider.Args{
		Name:           name,
		StartURL:       startURL,
		AllowedURL:     allowedURL,
		EntityURL:      entityURL,
		EntitySelector: entitySelector,
		UseBrowser:     useBrowser,
		Depth:          depth,
		MongoDbResource: "f/spider/resource/mongodb",
	}

	return name, spider.Start(args)
}
```

## URL Pattern Placeholders

In defining URL patterns for routing, filtering, or matching purposes, our system supports the use of specific placeholders within the pattern strings. These placeholders allow for dynamic matching against various URL structures. When a URL pattern includes one or more of these placeholders, it is transformed into a regular expression that matches the corresponding URL structure.

### Available Placeholders

- `{dir}`: Matches any sequence of characters except for slashes (`/`). Used to represent directory names in a URL path.
- `{any}`: Matches any sequence of characters, including an empty sequence. This is the most flexible placeholder and can match any part of a URL.
- `{some}`: Matches any non-empty sequence of characters. Similar to `{any}`, but requires at least one character to be present.
- `{num}`: Matches any sequence of digits. Useful for matching numeric identifiers within URLs.

### Placeholder Usage

Placeholders can be inserted into URL patterns to specify the parts of the URL that can vary. For example:

- `https://example.com/articles/{dir}/{some}`: Matches URLs that follow the structure of two path segments following `/articles/`, where the first segment can be any directory name, and the second must be a non-empty sequence.
- `https://example.com/products/{num}/details`: Matches URLs that include a numeric product ID followed by `/details`.

### Patterns Without Placeholders

If a URL pattern does not contain any of the specified placeholders, the system interprets the pattern in two ways:

1. **As a Regular Expression (If Not Empty):** If the pattern is not empty and contains no placeholders, it is treated as a ready-made regular expression. This allows for advanced matching scenarios where specific regular expression features are needed. When defining such patterns, ensure they are correctly escaped and adhere to regular expression syntax rules.

2. **Literal Match (If Empty):** An empty pattern matches nothing. This can be used to effectively disable a particular matching rule or filter.

### Escaping Special Characters

When placeholders are not used, and the pattern is intended as a regular expression, special characters must be escaped to prevent them from being interpreted as regex operators. For instance, `.` should be written as `\\.` and `/` as `\\/` to match these characters literally in URLs.

This flexible system of placeholders and direct regular expression input allows for precise control over URL matching, accommodating a wide range of routing and filtering requirements.




## Параметры

### Стартовый адрес
Первый адрес, который посетит парсер. Отсюда начинается обход страниц. Обычно это страница, котороя содержит наибольшее число ссылок на целевые страницы. В случае парсинга всего сайта - это главная страница или страница каталога. Для тематических выборок подойдёт страница оглавления категории.
В общем случае - эта та страница на сайте, которая поможет как можно скорее выявить все целевые страницы.

### Разрешённые адреса
Это адреса, на которых парсер будет искать ссылки на целевые страницы. Это могут быть страницы каталога, страницы категорий, страницы тегов, страницы поиска и т.д. В общем случае - это страницы, которые содержат ссылки на целевые страницы.
Параметр помогает роботу не уходить в сторону от целевых страниц. Например, обходить страницы административной части сайта, страницы с контактами и т.д.

### Целевые адреса
Это адреса, на которых размещены данные для извлечения. Это могут быть страницы товаров, страницы статей, страницы новостей и т.д. В общем случае - это страницы, которые содержат данные, которые нужно извлечь.

## Правила

Целевые страницы должны быть в диапазоне разрешённых адресов. Парсер не будет обходить страницы, которые не попадают под правила разрешённых адресов.
Стартовый адрес должен быть в диапазоне разрешённых адресов. Парсер не начнёт обход, если стартовый адрес не попадает под правила разрешённых адресов.
По умолчанию, разрешённые адреса включают в себя стартовый адрес и все его поддиректории. Таким образом, если стартовый адрес - это главная страница сайта, то парсер будет обходить все страницы сайта.
А если категории, то парсер будет обходить все страницы категории и подкатегорий, если адресное пространство сайта носит иерахический характер.


## Шаблоны

### Завершение с `{any}` и `{some}`

Параметры завершения адреса `{any}` и `{some}` позволяют уточнить, что именно должно быть в конце адреса.
- `{any}` позволяет указать, что в конце адреса может быть что угодно, включая пустую строку.
- `{some}` позволяет указать, что в конце адреса должно быть хотя бы один символ.

#### Обратите внимание
Шаблон `http://bbc.com/news{some}` не позволит парсеру обойти страницу `http://bbc.com/news`, так как в конце адреса должен быть хотя бы один символ.

### Вставка `{dir}`

Параметр `{dir}` позволяет указать, что в адресе должно быть название директории. Он подменяет любые слова между слешами `/`.
В отличие от `{any}` и `{some}`, парметр `{dir}` может быть использован внутри шаблона, а не только в конце.
Например, `http://bbc.com/news/{dir}/article/{some}` позволит парсеру обойти страницы `http://bbc.com/news/world/article/merry-xmas`, `http://bbc.com/news/sport/article/happy-ny-2024` и т.д.

### Вставка `{one,two,three}`
Параметр похож на `{dir}`, но позволяет указать несколько вариантов. В указанном месте адреса должно быть одно из указанных слов.
Например, `http://bbc.com/{news,sport}/article/{some}` позволит парсеру обойти страницы `http://bbc.com/news/article/happy-ny-2024`, `http://bbc.com/sport/article/happy-ny-2024` и т.д.

### Вставка `{num}`
Часто для добавления уникальности адреса могут содержать идентификаторы, которые являются числами.
Параметр `{num}` позволяет указать, что в адресе должно быть число. Он подменяет любые цифры.

Например, `http://bbc.com/news/article/news-number-{num}` позволит парсеру обойти страницы `http://bbc.com/news/article/news-number-1`, `http://bbc.com/news/article/news-number-2` и т.д.
Но пропустить страницу `http://bbc.com/news/article/news-number-` или `http://bbc.com/news/article/news-number-abc`.


## Пример

- Стартовый адрес: `http://bbc.com/news`
- Разрешённые адреса: `http://bbc.com/news{any}`
- Целевые адреса: `http://bbc.com/news/{dir}/article/{some}`

### Обяснение

- Парсер выберет все ссылки на странице `http://bbc.com/news`, которые начинаются с `http://bbc.com/news`.
- Посетит все страницы и вновь выберет удовлетворяющие условиям ссылки.
- Если одна из таких ссылок будет `http://bbc.com/news/world/article/happy-ny-2024`, то парсер извлечёт данные с этой страницы.
  В этом примере, `{dir}` подставится в `world`, а `{some}` в `happy-ny-2024`.


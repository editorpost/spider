
### Структура `Article`

Структура для сущности `Article`, которая будет достаточно универсальной и гибкой, чтобы покрывать большинство случаев.
Важно сохранить баланс между простотой и универсальностью, чтобы структура могла быть использована для различных типов контента и медиа.

1. **ID** (string) - Уникальный идентификатор статьи.
2. **Title** (string) - Заголовок статьи.
3. **Byline** (string) - Автор или авторы статьи.
4. **Content** (string) - Основной контент статьи в HTML формате.
5. **TextContent** (string) - Основной текст статьи без HTML тегов.
6. **Excerpt** (string) - Краткое содержание или анонс статьи.
7. **Images** ([]Image) - Массив изображений, связанных со статьей.
8. **Videos** ([]Video) - Массив видео, связанных со статьей.
9. **Quotes** ([]Quote) - Массив цитат из социальных сетей.
10. **PublishDate** (time.Time) - Дата публикации статьи.
11. **ModifiedDate** (time.Time) - Дата последнего изменения статьи.
12. **Tags** ([]string) - Теги, связанные со статьей.
13. **Source** (string) - Источник статьи (URL).
14. **Language** (string) - Язык статьи.
15. **Category** (string) - Категория статьи.
16. **SiteName** (string) - Название сайта, с которого была взята статья.
17. **AuthorSocialProfiles** ([]SocialProfile) - Массив социальных профилей автора.

### Структуры для вложенных элементов

#### Image

```go
type Image struct {
    URL     string `json:"url"`
    AltText string `json:"alt_text"`
    Width   int    `json:"width,omitempty"`
    Height  int    `json:"height,omitempty"`
    Caption string `json:"caption,omitempty"`
}
```

#### Video

```go
type Video struct {
    URL     string `json:"url"`
    EmbedCode string `json:"embed_code,omitempty"`
    Caption string `json:"caption,omitempty"`
}
```

#### Quote

```go
type Quote struct {
    Text     string `json:"text"`
    Author   string `json:"author"`
    Source   string `json:"source"`
    Platform string `json:"platform"` // e.g., Twitter, Facebook, Instagram
}
```

#### SocialProfile

```go
type SocialProfile struct {
    Platform string `json:"platform"` // e.g., Twitter, Facebook, Instagram
    URL      string `json:"url"`
}
```

### Полная структура `Article`

```go
type Article struct {
    ID                string          `json:"id"`
    Title             string          `json:"title"`
    Byline            string          `json:"byline"`
    Content           string          `json:"content"`
    TextContent       string          `json:"text_content"`
    Excerpt           string          `json:"excerpt"`
    Images            []Image         `json:"images"`
    Videos            []Video         `json:"videos"`
    Quotes            []Quote         `json:"quotes"`
    PublishDate       time.Time       `json:"publish_date"`
    ModifiedDate      time.Time       `json:"modified_date"`
    Tags              []string        `json:"tags"`
    Source            string          `json:"source"`
    Language          string          `json:"language"`
    Category          string          `json:"category"`
    SiteName          string          `json:"site_name"`
    AuthorSocialProfiles []SocialProfile `json:"author_social_profiles"`
}
```

### Обоснование выбора полей

1. **ID**: Уникальный идентификатор для каждой статьи.
2. **Title**: Заголовок статьи, который может быть использован для отображения и SEO.
3. **Byline**: Информация об авторе или авторах.
4. **Content**: Основной контент статьи в HTML формате, что позволяет включать различные виды медиа и разметки.
5. **TextContent**: Чистый текст статьи, полезный для анализа текста.
6. **Excerpt**: Краткий анонс статьи для отображения в списках или на главной странице.
7. **Images**: Массив изображений, связанных со статьей, с дополнительной информацией (альтернативный текст, размеры, подписи).
8. **Videos**: Массив видео с возможностью вставки embed-кодов и добавления подписей.
9. **Quotes**: Массив цитат из социальных сетей, что позволяет легко добавлять контент из Twitter, Facebook и Instagram.
10. **PublishDate**: Дата публикации статьи для сортировки и отображения.
11. **ModifiedDate**: Дата последнего изменения статьи, полезна для отслеживания обновлений.
12. **Tags**: Теги для классификации и поиска статей.
13. **Source**: URL источника статьи для отслеживания оригинального контента.
14. **Language**: Язык статьи для локализации.
15. **Category**: Категория статьи для структурирования контента.
16. **SiteName**: Название сайта, с которого была взята статья, для отображения и анализа источников.
17. **AuthorSocialProfiles**: Социальные профили автора для отображения дополнительных данных о авторе.

Эта структура позволяет хранить все необходимые данные для создания и модификации статей, а также для включения медиа и социального контента, что делает её универсальной и гибкой для различных видов публикаций.
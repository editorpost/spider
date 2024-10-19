### `collect` Package Documentation

#### Configuration

The `collect` package uses a configuration struct to manage its settings:

```go
type Config struct {
    StartURL        string `json:"StartURL"`
    AllowedURL      string `json:"AllowedURL"`
    ExtractURL      string `json:"ExtractURL"`
    ExtractSelector string `json:"ExtractSelector"`
    ExtractLimit    int    `json:"ExtractLimit"`
    UseBrowser      bool   `json:"UseBrowser"`
    Depth           int    `json:"Depth"`
    UserAgent       string `json:"UserAgent"`
    ProxyEnabled    bool   `json:"ProxyEnabled"`
    ProxySources    []string `json:"ProxySources"`
}
```

- **StartURL**: The initial URL to start scraping.
- **AllowedURL**: Regex to match URLs to reduce the number of URLs visited.
- **ExtractURL**: Regex to match entity URLs for extraction.
- **ExtractSelector**: CSS selector for extracting entities and filtering pages (default is `html`).
- **ExtractLimit**: Limit of entities to extract before stopping.
- **UseBrowser**: Flag to use a browser for rendering the page.
- **Depth**: Depth for link following; `1` means only links on the scraped page are visited.
- **UserAgent**: User agent string for the collector.
- **ProxyEnabled**: Flag to enable proxy usage.
- **ProxySources**: List of proxy sources.

#### Architecture

The `collect` package is built around the `colly` library for web scraping. Key components include:

- **Dispatcher**: Sets up event handlers for HTML elements, errors, requests, and responses.
- **Crawler**: Main struct that manages the scraping process, including initializing the collector and handling proxies.
- **Browser Interface**: Supports using a headless browser for scraping JavaScript-rendered content.

#### Mechanism of Collection

The collection mechanism involves setting up a `colly.Collector` with various handlers:

1. **Initialization**: The `NewDispatcher` function initializes the dispatcher with configuration and dependencies.
2. **Event Handlers**: The `WithDispatcher` function sets up handlers for different events (HTML elements, errors, requests, responses, scraped data).
3. **Queue Management**: A request queue is used to manage URLs to be visited.

```go
func (crawler *Dispatch) request(r *colly.Request) {
    crawler.deps.Monitor.OnRequest(r)
}

func (crawler *Dispatch) response(r *colly.Response) {
    crawler.deps.Monitor.OnResponse(r)
}

func (crawler *Dispatch) scraped(r *colly.Response) {
    crawler.deps.Monitor.OnScraped(r)
}
```

#### Diagram of Events Relation

Here is a simplified diagram of event relations:

1. **Collector Initialization**:
    - `collector := colly.NewCollector(...)`

2. **Event Handlers Setup**:
    - `collector.OnHTML('a[href]', visitHandler)`
    - `collector.OnHTML('html', extractHandler)`
    - `collector.OnError(errorHandler)`
    - `collector.OnRequest(requestHandler)`
    - `collector.OnResponse(responseHandler)`
    - `collector.OnScraped(scrapedHandler)`

3. **Queue Management**:
    - `queue := queue.New(...)`
    - `queue.AddURL(startURL)`

4. **Scraping Process**:
    - `collector.Visit(startURL)`

This documentation should help you understand and configure the `collect` package in the `editorpost/spider` repository.
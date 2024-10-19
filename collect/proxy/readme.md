### `proxy` Package Documentation

#### Principles

The `proxy` package is designed to manage proxy configurations and rotations for the scraping tasks. It allows the `colly` collector to use a pool of proxies to distribute requests, improving anonymity and avoiding IP bans.

#### Limitations and Caveats

- **Proxy Availability**: The effectiveness of the proxy pool depends on the availability and reliability of the proxies. If proxies are down or slow, the scraping tasks may be affected.
- **Configuration**: Proper configuration of proxy sources is essential. Incorrect or missing configurations can lead to failures in the proxy setup.
- **Performance**: The performance may vary depending on the quality and speed of the proxies used.

#### Error Handling Steps and Retries

The `proxy` package includes mechanisms to handle errors and retries:

1. **Error Handling**:
    - If a proxy fails, it is logged and the request is retried with a different proxy.
    - Errors in setting up proxies are returned immediately, preventing the collector from starting with invalid configurations.

2. **Retries**:
    - The package uses retries for both proxy-related errors and response errors.
    - The retry mechanism is configurable, allowing the number of retries and the conditions under which retries are performed to be adjusted.

#### Configuration Example
```go
const (
    DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
)

type Config struct {
    StartURL     string   `json:"StartURL"`
    ProxyEnabled bool     `json:"ProxyEnabled"`
    ProxySources []string `json:"ProxySources"`
}
```

#### Setting Up Proxies
```go
func WithProxyPool(args *config.Config) (colly.CollectorOption, error) {
    var (
        err       error
        poolReady bool
        proxies   *proxy.Pool
    )

    if args.ProxyEnabled {
        proxies, err = proxy.StartPool(args.StartURL, args.ProxySources...)
        if err != nil {
            return nil, err
        }
        poolReady = true
    }

    return func(c *colly.Collector) {
        if poolReady {
            c.WithTransport(proxies.Transport())
        }
        c.SetRequestTimeout(15 * time.Second)
    }, nil
}
```

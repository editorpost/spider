package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// LoadJSONList loads the valid list from the given url
// Returns nil if the url is empty.
func LoadJSONList(url string) []string {

	if url == "" {
		return nil
	}

	// fetch the url
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// parse the response body
	var proxies []*Proxy
	if err = json.Unmarshal(body, &proxies); err != nil {
		panic(err)
	}

	var lines []string
	for _, p := range proxies {
		lines = append(lines, p.String())
	}

	if len(lines) == 0 {
		panic("no proxies found")
	}

	return lines
}

// LoadStringList loads the valid list from proxyscrape.com
func LoadStringList(sourceURL string) ([]string, error) {

	// fetch the url
	res, err := http.Get(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("can not load proxy list from proxyscrape.com: %w", err)
	}
	defer res.Body.Close()

	// parse the response body
	lines := []string{}
	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// LoadStringLists loads the valid list from public sources
func LoadStringLists(sources []string) ([]string, error) {

	if len(sources) == 0 {
		return nil, nil
	}

	wg := &sync.WaitGroup{}

	// load all sources
	proxies := NewList()

	for _, source := range sources {
		wg.Add(1)

		go func(source string, wg *sync.WaitGroup) {
			defer wg.Done()
			urls, err := LoadStringList(source)
			if err != nil {
				return
			}

			proxies.Add(NewProxies(urls...)...)
		}(source, wg)
	}

	wg.Wait()

	return proxies.Strings(), nil
}

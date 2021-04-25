// see https://tour.golang.org/concurrency/10
// Exercise: Web Crawler
// In this exercise you'll use Go's concurrency features to parallelize a web crawler.
// Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.
// Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!

package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeCache struct {
	urls map[string]struct{}
	mux sync.Mutex
}

func (sc *SafeCache) Has(url string) bool {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	_, ok := sc.urls[url]
	return ok
}

func (sc *SafeCache) Push(url string) {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	sc.urls[url] = struct{}{}
}

type Crawler struct {
	cache *SafeCache
	fetcher Fetcher
}

func (crl *Crawler) Crawl (url string, depth int) {
	if depth <= 0 || crl.cache.Has(url) {
		return
	}
	crl.cache.Push(url)
	body, urls, err := crl.fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		go crl.Crawl(u, depth-1)
	}
}

func main() {
	crl := Crawler{
		cache: &SafeCache{urls: make(map[string]struct{})}, 
		fetcher: fetcher,
	}
	go crl.Crawl("https://golang.org/", 4)
	time.Sleep(time.Second)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

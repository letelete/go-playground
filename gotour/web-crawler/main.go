// see https://tour.golang.org/concurrency/10
// Exercise: Web Crawler
// In this exercise you'll use Go's concurrency features to parallelize a web crawler.
// Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.
// Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!

package main

import (
	"fmt"
	"sync"
)

// -- Fetcher --
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// -- CrawlResult --
type CrawlResult struct {
	url  string
	body string
	err  error
}

func (cr *CrawlResult) String() string {
	if cr.err != nil {
		return fmt.Sprint("error: ", cr.url, " message: ", cr.err)
	}
	return fmt.Sprint("found: ", cr.url, " - ", cr.body)
}

// -- SafeCache --
type SafeCache struct {
	urls map[string]struct{}
	mux  sync.Mutex
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

// -- Crawler --
type Crawler struct {
	cache   *SafeCache
	fetcher Fetcher
}

func (crl *Crawler) Crawl(url string, depth int, cr chan<- *CrawlResult) {
	defer close(cr)
	wg := sync.WaitGroup{}
	wg.Add(1)
	crl.crawlRecursively(url, depth, cr, &wg)
	wg.Wait()
}

func (crl *Crawler) crawlRecursively(url string, depth int, cr chan<- *CrawlResult, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 || crl.cache.Has(url) {
		return
	}
	crl.cache.Push(url)
	body, urls, err := crl.fetcher.Fetch(url)
	cr <- &CrawlResult{url: url, err: err, body: body}
	if err != nil {
		return
	}
	for _, u := range urls {
		wg.Add(1)
		go crl.crawlRecursively(u, depth-1, cr, wg)
	}
}

// -- Main --

func printCrawlResults(cr chan *CrawlResult) {
	for {
		v, ok := <-cr
		if !ok {
			break
		}
		fmt.Println(v)
	}
}

func main() {
	crl := Crawler{
		cache:   &SafeCache{urls: make(map[string]struct{})},
		fetcher: fetcher,
	}
	cr := make(chan *CrawlResult)
	go crl.Crawl("https://golang.org/", 4, cr)
	printCrawlResults(cr)
}

// -- Fake Fetcher Implementation --

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

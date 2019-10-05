package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeMap struct {
	m   map[string]bool
	mux sync.Mutex
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	safemap := SafeMap{m: make(map[string]bool)}
	msgCh := make(chan string)	// to show msg from subprocesses
	endCh := make(chan bool)	// to wait for all subprocesses
	go doCrawl(url, depth, fetcher, &safemap, msgCh, endCh)
	fmt.Printf("Result: %s\n", <-msgCh)
	<-endCh	
}

func doCrawl(url string, depth int, fetcher Fetcher, safemap *SafeMap, msgCh chan string, endCh chan bool) {
	if depth <= 0 {
		msgCh <- fmt.Sprintf("url %s depth <= 0", url)
		endCh <- false
		return
	}
	
	// just check whether this url has been crawled or not
	// no need to lock
	if _, ok := safemap.m[url]; ok {
		msgCh <- fmt.Sprintf("url %s crawled before", url)
		endCh <- false
		return
	}
	
	fmt.Printf("Begin Task: %s\n", url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		msgCh <- fmt.Sprintf("url %s in error %s", url, err)
		// we avoid crawl same error url next time
		safemap.mux.Lock()
		safemap.m[url] = false
		safemap.mux.Unlock()
		endCh <- false
		return
	}
	
	// lock for writting 
	safemap.mux.Lock()
	
	if _, ok := safemap.m[url]; ok {
		// crawled before
		safemap.mux.Unlock()
		msgCh <- fmt.Sprintf("url %s crawled before", url)
		endCh <- false
		return
	}

	safemap.m[url] = true // record to safemap
	safemap.mux.Unlock()

	// show new results
	fmt.Printf("found: %s %q\n", url, body)
	msgCh <- fmt.Sprintf("url %s finished", url)
	
	newMsgCh := make(chan string) // to show msg from subprocesses
	newEndCh := make(chan bool) // to wait for all subprocesses
	for _, u := range urls {
		go doCrawl(u, depth-1, fetcher, safemap, newMsgCh, newEndCh)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Printf("Result: %s\n", <-newMsgCh)
	}

	for i := 0; i < len(urls); i++ {
		<-newEndCh
	}
	endCh <- true // end itself
	return
}

// below is the original code without modification

func main() {
	Crawl("https://golang.org/", 4, fetcher)
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

package main

import (
	"fmt"
	"sync"
)

//-------------------------------------------------------------------------

func crawler(base, current string, visited *[]string) error {
	*visited = append(*visited, normalizeURL(current))
	html, err := getHTML(current)
	if err != nil {
		if err.Error() == "error Header: not text/html" {
			return nil
		}
		return err
	}
	urls, err := getURLsFromHTML(html, base)
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		return nil
	}
	for _, url := range urls {
		if !link_visited(url, *visited) {
			fmt.Println("//-----------")
			fmt.Printf("%s-> %s\n", current, url)
			fmt.Println("//-----------")
			err := crawler(base, url, visited)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//-------------------------------------------------------------------------

func link_visited(visiting string, visited []string) bool {
	normal := normalizeURL(visiting)
	for _, link := range visited {
		if link == normal {
			return true
		}
	}
	return false
}

//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+//

func start_crawler(arg string) error {
	cfg := get_concurrent_struct(arg)
	err := cfg.crawler_concurrent(arg)
	if err != nil {
		return err
	}
	cfg.wg.Wait()
	close(cfg.concurrencyControl)
	return nil
}

//-------------------------------------------------------------------------

type Config struct {
	visit            map[string]bool
	baseURL            string
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

//-------------------------------------------------------------------------

func get_concurrent_struct(raw_base_URL string) *Config {
	return &Config{
		visit: map[string]bool{},
		baseURL: raw_base_URL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 2),
		wg: &sync.WaitGroup{},
	}
}

//-------------------------------------------------------------------------

func (cfg *Config) crawler_concurrent(CurrentURL string) error {
	normal := normalizeURL(CurrentURL)
	cfg.mu.Lock()
	cfg.visit[normal] = true
	cfg.mu.Unlock()
	fmt.Printf("-> %s\n", CurrentURL)
	//*************
	html, err := getHTML(CurrentURL)
	if err != nil {
		if err.Error() == "error Header: not text/html" {
			return nil
		}
		return err
	}
	//*************
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		return err
	}
	//*************
	if len(urls) == 0 {
		return nil
	}
	//*************
	cfg.mu.Lock()
	cfg.update_visit_concurrent(urls)
	cfg.mu.Unlock()

	for _, url := range urls {
		cfg.wg.Add(1)
		go func(url string) {
			cfg.concurrent_call(url)
		}(url)
	}
	return nil
}

//-------------------------------------------------------------------------

func (cfg *Config) concurrent_call(url string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	if !cfg.visit[normalizeURL(url)] {
		err := cfg.crawler_concurrent(url)
		if err != nil {
			fmt.Println(err)
		}
	}
}


//-------------------------------------------------------------------------

func (cfg *Config) update_visit_concurrent(to_visit []string) {
	for _, url := range to_visit {
		normal := normalizeURL(url)
		if _, ok := cfg.visit[normal]; !ok {	
			cfg.visit[normal] = false
		}
	}
}

//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+//

func get_urls(CurrentURL string) {
	fmt.Println("get_urls...")
	html, err := getHTML(CurrentURL)
	if err != nil {
		fmt.Println(err)
	}
	//*************
	urls, err := getURLsFromHTML(html, CurrentURL)
	if err != nil {
		fmt.Println(err)
	}
	for _, link := range urls {
		fmt.Println(link)
	}
}
package main

import (
	"fmt"
	"sync"
	"errors"
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

func start_crawler(base string, concurrency, max_visit int) error {
	cfg := get_concurrent_struct(base, concurrency, max_visit)
	err := cfg.crawler_concurrent(base)
	if err != nil {
		return err
	}
	cfg.wg.Wait()
	cfg.print_out()
	close(cfg.concurrencyControl)
	return nil
}

//-------------------------------------------------------------------------

func (cfg Config) print_out() {
	fmt.Println("////////////////////////////")
	fmt.Printf("REPORT for %s\n", cfg.baseURL)
	fmt.Println("////////////////////////////")
	fmt.Println("")
	for key, list := range cfg.visit {
		fmt.Println("    =========================")
		fmt.Printf("    REPORT for %s\n", key)
		fmt.Println("    =========================")
		fmt.Println("    Found internal links to:")
		fmt.Println("    -------------------------")
		for _, item := range list {
			fmt.Println("        "+item)
		}
		fmt.Println("")
	}
}

//-------------------------------------------------------------------------

type Config struct {
	exit 				bool
	max_visit			int
	visit				map[string][]string
	baseURL				string
	mu					*sync.Mutex
	concurrencyControl	chan struct{}
	wg					*sync.WaitGroup
}

//-------------------------------------------------------------------------

func get_concurrent_struct(raw_base_URL string, concurrency, max_visit int) *Config {
	return &Config{
		exit: 				false,
		max_visit: 			max_visit,
		visit: 				map[string][]string{},
		baseURL: 			raw_base_URL,
		mu: 				&sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrency),
		wg: 				&sync.WaitGroup{},
	}
}

//-------------------------------------------------------------------------

func (cfg *Config) crawler_concurrent(CurrentURL string) error {
	fmt.Printf("-> %s\n", CurrentURL)
	//*************
	cfg.mu.Lock()
	urls, err := cfg.get_urls(CurrentURL)
	cfg.mu.Unlock()
	if err != nil {
		cfg.mu.Lock()
		cfg.exit = true
		cfg.mu.Unlock()
		return err
	} else if len(urls) == 0 {
		return nil
	} else if cfg.exit {
		return errors.New("Max visit reached!!!")
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
		if cfg.exit {
			return nil
		}
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
	if cfg.exit {
		return
	}
	if len(cfg.visit[normalizeURL(url)]) == 0 {
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
			cfg.visit[normal] = []string{}
		}
	}
}

//-------------------------------------------------------------------------

func (cfg *Config) reach_max_visit() {
	i := 0
	for _, visited := range cfg.visit {
		if len(visited) != 0 {
			i++
		}
	}
	if i >= cfg.max_visit {
		cfg.exit = true
	}
}

//-------------------------------------------------------------------------

func (cfg *Config) get_urls(CurrentURL string) ([]string, error) {
	html, err := getHTML(CurrentURL)
	if err != nil {
		if error.Error(err) == "error Header: not text/html" {
			return []string{}, nil
		}
		return []string{}, err
	}
	//*************
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		return []string{}, err
	}
	//*************
	cfg.visit[normalizeURL(CurrentURL)] = urls
	cfg.reach_max_visit()
	return urls, nil
}

//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+//

func get_urls(CurrentURL string) {
	html, _ := getHTML(CurrentURL)
	//*************
	urls, _ := getURLsFromHTML(html, CurrentURL)
	for _, url := range urls {
		fmt.Println(url)
	}
}
package main

import (
	"fmt"
)

//-------------------------------------------------------------------------

func crawler(base, current string, visited *[]string) error{
	normal, err := normalizeURL(current)
	if err != nil {
		return err
	}
	*visited = append(*visited, normal)
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
	normal, _ := normalizeURL(visiting)
	for _, link := range visited {
		if link == normal {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------------
package main

import (
	"net/url"
	"golang.org/x/net/html"
	"fmt"
	"strings"
)

//-------------------------------------------------------------------------

func normalizeURL(raw_url string) string {
	prefixes := []string{"http:/", "https:/", "ftp:/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(raw_url, prefix) && !strings.HasPrefix(raw_url, prefix+"/") {
			raw_url = strings.Replace(raw_url, prefix, prefix+"/", 1)
			break  
		}
	}
	for strings.Contains(raw_url, "///") {
		raw_url = strings.Replace(raw_url, "///", "//", -1)
	}
	url_struct, err := url.Parse(raw_url)
	if err != nil {
		fmt.Println(err)
	}
	str_h := strings.Trim(url_struct.Host, "/")
	str_p := strings.Trim(url_struct.Path, "/")
	str := str_h + "/" + str_p+"/"
	str = strings.Replace(str, "//", "/", -1)
	return str
}

//-------------------------------------------------------------------------

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if string(a.Val[0]) == "/" {
						links = append(links, rawBaseURL+a.Val)
					} else if string(a.Val[0]) == "#" {
						continue
					} else {
						links = append(links, a.Val)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	var final []string
	for i, link := range links {
		good := true
		for j:=i+1; j+1<len(links); j++ {
			if link == links[j] {
				good = false
			}
		}
		if good {
			final = append(final, link)
		}
	}
	return final, nil
}

//-------------------------------------------------------------------------
package main

import (
	"net/url"
	"golang.org/x/net/html"
	"fmt"
	"os"
	"sort"
	"regexp"
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
	
	str := str_h
	if str_p != "" { str += "/" + str_p }
	
	str = strings.Replace(str, "//", "/", -1)
	return str
}

//-------------------------------------------------------------------------

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	siteRoot := base.Scheme + "://" + base.Host 
	seen := make(map[string]bool)
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				href := strings.TrimSpace(a.Val)
				if href == "" || strings.HasPrefix(href, "#") {
					continue
				}

				// Resolve relative URLs correctly
				ref, err := url.Parse(href)
				if err != nil {
					continue
				}
				absolute := base.ResolveReference(ref)

				full := absolute.String()

				// Only keep links that belong to ingatlan.com
				if !strings.HasPrefix(full, siteRoot) {
					continue
				}

				// Only keep links that end with a pure number (listing ID)
				// Examples: .../35503400  or  .../35503400/
				path := strings.TrimSuffix(absolute.Path, "/")
				parts := strings.Split(path, "/")
				last := parts[len(parts)-1]

				if matched, _ := regexp.MatchString(`^\d+$`, last); matched {
					if !seen[full] {
						seen[full] = true
						links = append(links, full)
					}
				}
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	var final []string
	base_url_struct, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, err
	}
	normal_base := normalizeURL(base_url_struct.Host)
	for i, link := range links {
		current_url_struct, err := url.Parse(link)
		if err != nil {
			return []string{}, err
		}

		normal_current := normalizeURL(current_url_struct.Host)
		if normal_base != normal_current {
			continue
		}

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

func saveLinksToFile(filename string, newLinks []string) error {
	// Load previous links (if the file exists)
	oldLinks := make(map[string]bool)
	data, err := os.ReadFile(filename)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				oldLinks[line] = true
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading old file: %w", err)
	}

	// Build set of new links
	newSet := make(map[string]bool)
	for _, link := range newLinks {
		link = strings.TrimSpace(link)
		if link != "" {
			newSet[link] = true
		}
	}

	// Detect changes
	var added, removed []string

	for link := range newSet {
		if !oldLinks[link] {
			added = append(added, link)
		}
	}
	for link := range oldLinks {
		if !newSet[link] {
			removed = append(removed, link)
		}
	}

	// Print warnings
	if len(added) > 0 {
		fmt.Printf("⚠️  %d new link(s) added:\n", len(added))
		for _, l := range added {
			fmt.Println("  +", l)
		}
	}
	if len(removed) > 0 {
		fmt.Printf("⚠️  %d link(s) removed:\n", len(removed))
		for _, l := range removed {
			fmt.Println("  -", l)
		}
	}
	if len(added) == 0 && len(removed) == 0 {
		fmt.Println("No changes detected.")
	}

	// Write the new list (sorted for nicer diffs)
	sorted := make([]string, 0, len(newSet))
	for link := range newSet {
		sorted = append(sorted, link)
	}
	sort.Strings(sorted)

	content := strings.Join(sorted, "\n") + "\n"
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Saved %d links to %s\n", len(sorted), filename)
	return nil
}
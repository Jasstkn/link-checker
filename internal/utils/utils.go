// Package utils provides independent untils functions for the linkchecker package
package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	supportedSchemes = [...]string{"http://", "https://"}
)

// ParseHTML parses HTML and return []string of links
func ParseHTML(body string) []string {
	re := regexp.MustCompile(`<a href="(http[a-zA-Z\d\-_.:\/]*?)"`)
	matched := re.FindAllStringSubmatch(body, -1)

	links := make([]string, 0, len(matched))
	for _, v := range matched {
		links = append(links, v[1])
	}
	return links
}

// ValidateLinks validates []string and return number of broken with list of them
func ValidateLinks(links []string) (int, []string) {
	var wg sync.WaitGroup
	ch := make(chan string)

	for _, l := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			client := &http.Client{}
			req, _ := http.NewRequest("GET", l, nil)
			resp, _ := client.Do(req)
			// linkedin return 999
			if resp.StatusCode > 400 && resp.StatusCode != 999 {
				ch <- l
			}
		}(l)
	}

	// close chanel to prevent memory leak
	go func() {
		// wait until counter is 0
		wg.Wait()
		close(ch)
	}()

	var brokenLinks []string
	for l := range ch {
		brokenLinks = append(brokenLinks, l)
	}

	return len(brokenLinks), brokenLinks
}

// ValidateURL returns whether a given URL scheme is supported
func ValidateURL(url string) error {
	for _, scheme := range supportedSchemes {
		if strings.Contains(url, scheme) {
			return nil
		}
	}
	return fmt.Errorf("missing or not supported URL scheme in %q. Available: %s",
		url, strings.Join(supportedSchemes[:], ", "))
}

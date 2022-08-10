package utils

import (
	"net/http"
	"regexp"
	"sync"
)

func ParseHtml(body string) []string {
	re := regexp.MustCompile(`<a href="(http[a-zA-Z\d\-_.:\/]*?)"`)
	matched := re.FindAllStringSubmatch(body, -1)

	links := make([]string, 0, len(matched))
	for _, v := range matched {
		links = append(links, v[1])
	}
	return links
}

func ValidateLinks(links []string) (int, []string) {
	var wg sync.WaitGroup
	ch := make(chan string)

	for _, l := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			r, _ := http.Get(l)
			if r.StatusCode > 299 {
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

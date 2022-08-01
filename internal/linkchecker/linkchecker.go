package linkchecker

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

func ParseHtml(body string) []string {
	re := regexp.MustCompile(`<a href="(http.*?)"`)
	matched := re.FindAllStringSubmatch(body, -1)

	links := make([]string, 0, len(matched))
	for _, v := range matched {
		links = append(links, v[1])
	}
	return links
}

func ValidateLinks(links []string) (n int, brokenLinks []string) {
	var wg sync.WaitGroup
	for _, l := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			_, err := http.Get(l)
			if err != nil {
				n++
				brokenLinks = append(brokenLinks, l)
			}
		}(l)
	}

	go func() {
		// wait until counter is 0
		wg.Wait()
	}()
	return n, brokenLinks
}

func LinkChecker(url string) (string, error) {
	req, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	links := ParseHtml(string(body))

	if len(links) == 0 {
		return "No links were found", nil
	}

	brokenNum, brokenLinks := ValidateLinks(links)

	if brokenNum == 0 {
		return fmt.Sprintf("%d links scanned, %d broken found", len(links), brokenNum), nil
	}

	return fmt.Sprintf("%d links scanned, %d broken links found, %s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
}

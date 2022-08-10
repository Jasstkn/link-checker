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

	switch {
	case len(links) == 1 && brokenNum == 0:
		return fmt.Sprintf("%d link scanned, %d broken links found%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	case len(links) == 1 && brokenNum == 1:
		return fmt.Sprintf("%d link scanned, %d broken link found:\n%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	case brokenNum == 1:
		return fmt.Sprintf("%d links scanned, %d broken link found:\n%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	default:
		return fmt.Sprintf("%d links scanned, %d broken links found:\n%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	}

	// switch brokenNum {
	// case 0:
	// 	return fmt.Sprintf("%d links scanned, %d broken links found%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	// case 1:
	// 	return fmt.Sprintf("%d links scanned, %d broken link found:\n%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	// default:
	// 	return fmt.Sprintf("%d links scanned, %d broken links found:\n%s", len(links), brokenNum, strings.Join(brokenLinks, ";\n")), nil
	// }
}

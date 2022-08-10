package linkchecker

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jasstkn/link-checker/internal/utils"
)

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
	links := utils.ParseHtml(string(body))

	if len(links) == 0 {
		return "No links were found", nil
	}

	brokenNum, brokenLinks := utils.ValidateLinks(links)

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
}

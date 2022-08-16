// Package to build a command-line tool for checking broken links
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Jasstkn/link-checker/pkg/linkchecker"
)

var (
	Version   string // Version stores release tag
	GitCommit string // GitCommit stores SHA's release commit

	supportedSchemes = [...]string{"http://", "https://"}
)

func main() {
	urlFlag := flag.String("url", "", "URL to check")
	versionFlag := flag.Bool("version", false, "Print the current version")
	flag.Parse()

	switch {
	case *versionFlag:
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		return
	}

	if err := ValidateURL(*urlFlag); err != nil {
		log.Fatal(err)
	}

	check, err := linkchecker.LinkChecker(*urlFlag)

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(check)
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

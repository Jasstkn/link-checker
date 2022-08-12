// Package to build a command-line tool for checking broken links
package main

import (
	"flag"
	"fmt"

	"github.com/Jasstkn/link-checker/pkg/linkchecker"
)

var (
	Version   string // Version stores release tag
	GitCommit string // GitCommit stores SHA's release commit
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

	check, err := linkchecker.LinkChecker(*urlFlag)

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(check)
}

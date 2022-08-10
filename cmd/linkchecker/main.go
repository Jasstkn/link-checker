package main

import (
	"flag"
	"fmt"

	"github.com/Jasstkn/link-checker/pkg/linkchecker"
)

var (
	Version   string
	GitCommit string
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

package main

import (
	"flag"
	"fmt"

	linkchecker "github.com/Jasstkn/link-checker/internal/link-checker"
)

func main() {
	url := flag.String("url", "", "url to check")
	flag.Parse()

	check, err := linkchecker.LinkChecker(*url)

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(check)
}

// Package to build a command-line tool for checking broken links
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Jasstkn/link-checker/internal/linkchecker"
	"github.com/Jasstkn/link-checker/internal/server"
	"github.com/Jasstkn/link-checker/internal/utils"
)

var (
	Version   string // Version stores release tag
	GitCommit string // GitCommit stores SHA's release commit
)

func main() {
	urlFlag := flag.String("url", "", "URL to check")
	versionFlag := flag.Bool("version", false, "Print the current version")
	serverFlag := flag.Bool("server", false, "Enable server mode")
	flag.Parse()

	switch {
	case *versionFlag:
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		return
	case *serverFlag && *urlFlag != "":
		log.Fatal("can't use CLI and server mode simultaneously")
	case *serverFlag:
		server.Init()
	case *urlFlag != "":
		break
	default:
		fmt.Println("run linkchecker -help to see available options")
		os.Exit(0)
	}

	if err := utils.ValidateURL(*urlFlag); err != nil {
		log.Fatal(err)
	}

	check, err := linkchecker.LinkChecker(*urlFlag)

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(check)
}

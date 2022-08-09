# A website link checker

[![Go Report Card](https://goreportcard.com/badge/github.com/Jasstkn/link-checker)](https://goreportcard.com/report/github.com/Jasstkn/link-checker)
[![codecov](https://codecov.io/gh/Jasstkn/link-checker/branch/master/graph/badge.svg?token=Q95TYDZXJP)](https://codecov.io/gh/Jasstkn/link-checker)

A linkchecker is a simple CLI tool to find all broken links in your website.

## Build

- with [Taskfile][1]

    ```bash
    task build
    ```

- with Go CLI

    ```bash
    go build -o linkcheker cmd/linkchecker/main.go
    ```

## Test


- with [Taskfile][1]

    ```bash
    task test
    # run all tests with integrational included
    task test-all
    ```

- with Go CLI

    ```bash
    go test ./...
    ```

## Example of usage

```bash
./linkcheker -url https://en.wikipedia.org/
parsed url: https://en.wikipedia.org/
97 links scanned, 0 broken links found

./linkcheker -url https://en.wikipedia-broken.org/
Get "https://en.wikipedia-broken.org/": dial tcp: lookup en.wikipedia-broken.org: no such host

./linkcheker -url https://github.com/Jasstkn/link-checker
13 links scanned, 1 broken link found:
https://github.com/Jasstkn/test-repo.git
```

[1]: https://taskfile.dev/

> broken link to test

<a href="https://github.com/Jasstkn/test-repo.git">test broken link</a>

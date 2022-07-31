# A website link checker

[![Go Report Card](https://goreportcard.com/badge/github.com/Jasstkn/link-checker)](https://goreportcard.com/report/github.com/Jasstkn/link-checker)

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

## Example of usage

```bash
./linkcheker -url https://en.wikipedia.org/
parsed url: https://en.wikipedia.org/
97 links scanned, 0 broken found

./linkcheker -url https://en.wikipedia-broken.org/
Get "https://en.wikipedia-broken.org/": dial tcp: lookup en.wikipedia-broken.org: no such host
```

[1]: https://taskfile.dev/

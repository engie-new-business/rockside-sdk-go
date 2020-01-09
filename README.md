[![GoDoc](https://godoc.org/github.com/rocksideio/rockside-sdk-go?status.svg)](https://godoc.org/github.com/rocksideio/rockside-sdk-go)

# Rockside GO SDK

Official Rockside SDK for GO applications.

## Usage

### Client library

```go

import "github.com/rocksideio/rockside-sdk-go"

...
client, err := rockside.NewClient(apiKey)
...
```

### Command Line Interface

Install it with

```sh
go get github.com/rocksideio/rockside-sdk-go/cmd/rockside
```

Then
```sh
rockside -h
```

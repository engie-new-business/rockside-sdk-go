[![GoDoc](https://godoc.org/github.com/rocksideio/rockside-sdk-go?status.svg)](https://godoc.org/github.com/rocksideio/rockside-sdk-go)

# Rockside GO SDK

Official Rockside SDK for GO applications.

## Client Library Usage

To use the client look at the [reference and examples](https://godoc.org/github.com/rocksideio/rockside-sdk-go) 

## CLI Usage

We provide a basic command line interface using the Rockside SDK GO, to interact with the Rockside API.

Install it with:

```sh
go get github.com/rocksideio/rockside-sdk-go/cmd/rockside
```

Display the various commands & flags available with:

```sh
rockside -h
```

Then to use commands export your API key:

```
export ROCKSIDE_API_KEY=...
rockside --tesnet --verbose identities ls
```
[![GoDoc](https://godoc.org/github.com/rocksideio/rockside-sdk-go?status.svg)](https://godoc.org/github.com/rocksideio/rockside-sdk-go)

# Rockside GO SDK

Official Rockside SDK for GO applications.

## Client Library Usage

To use the client look at the [reference and examples](https://godoc.org/github.com/rocksideio/rockside-sdk-go) 

## Command Line Interface Usage

To interact with the Rockside API, deploy contracts, etc. we provide a CLI (that uses the Rockside SDK GO)

### Install 

Get the latest CLI binary for [macOS, Windows or Linux here](https://github.com/rocksideio/rockside-sdk-go/releases)!

(... or if you have GO locally do: `go get github.com/rocksideio/rockside-sdk-go/cmd/rockside`)

### Usage

Display the various commands & flags available with:

```sh
rockside -h
```

Then to use commands export your API key:

```
export ROCKSIDE_API_KEY=...
rockside --tesnet --verbose identities ls
```

For instance you can deploy a contract with:

```
export ROCKSIDE_API_KEY=...
./rockside --testnet deploy-contract /tmp/mycontract.sol
```
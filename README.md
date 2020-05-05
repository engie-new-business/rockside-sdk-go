[![GoDoc](https://godoc.org/github.com/rocksideio/rockside-sdk-go?status.svg)](https://godoc.org/github.com/rocksideio/rockside-sdk-go)

# Rockside GO SDK

Official Rockside SDK for GO applications.

## Client Library Usage

To use the client look at the [reference and examples](https://pkg.go.dev/github.com/rocksideio/rockside-sdk-go?tab=overview) 

## Command Line Interface Usage

To interact with the Rockside API, deploy contracts, etc. we provide a CLI (that uses the Rockside SDK GO)

#### Install 

Get the latest CLI binary for [macOS, Windows or Linux here](https://github.com/rocksideio/rockside-sdk-go/releases)!

If you have GO locally do: `go get github.com/rocksideio/rockside-sdk-go/cmd/rockside`

To build from the cloned repo do: `go build ./cmd/rockside/; ./rockside -h`

#### Usage

Display the various commands & flags available with:

```sh
rockside -h
```

Then to use commands export your API key:

```sh
export ROCKSIDE_API_KEY=...
rockside --tesnet --verbose identities ls
```

For instance you can deploy a contract with:

```sh
export ROCKSIDE_API_KEY=...
rockside --testnet deploy-contract /tmp/mycontract.sol
```

Display/track a transaction:

```console
# Show a transaction using its hash
rockside transaction show 0x73da8b72acf620c05471edded3e425e853b6ad6853b8fcfd6adf754fff4bce9b

# Show a transaction using its tracking ID
rockside transaction show 01B7J50J5N7PEFMCY181N32938
```

Other useful commands:

```console
# List my identities
rockside --testnet identities ls

# Show a transaction receipt from a transaction hash
rockside --testnet receipt 0x97dfce42248a3f67f5a0660fab117b0ed7cb57af799bdda8854eca5ae5a98e28
```
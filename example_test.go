package rockside_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocksideio/rockside-sdk-go"
)

var client *rockside.Client

func ExampleNewClient() {
	client, err := rockside.NewClient("https://api.rockside.io", os.Getenv("ROCKSIDE_API_KEY"), rockside.Testnet)
	if err != nil {
		panic(err)
	}

	identities, err := client.Identities.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(identities)
}

func ExampleRPCClient() {
	// Get a RPC client from your existing Rockside client.
	rpc := client.RPCClient

	accounts, err := rpc.EthAccounts()
	if err != nil {
		panic(err)
	}
	fmt.Println(accounts)

	balance, err := rpc.BalanceAt(context.Background(), common.Address{}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
}

func ExampleDeployContract() {
	identities, err := client.Identities.List()
	if err != nil {
		panic(err)
	}

	var contractCode, jsonABI string
	txHash, err := client.DeployContractWithIdentity(identities[0], contractCode, jsonABI)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)
}

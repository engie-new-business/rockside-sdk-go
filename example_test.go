package rockside_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocksideio/rockside-sdk-go"
)

var rocksideClient *rockside.Client

func ExampleNewClientFromAPIKey() {
	rocksideAPIclient, err := rockside.NewClientFromAPIKey(os.Getenv("ROCKSIDE_API_KEY"), rockside.Testnet)
	if err != nil {
		panic(err)
	}

	smartWallets, err := rocksideAPIclient.SmartWallets.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(smartWallets)
}

func ExampleNewClientFromToken() {
	rocksideAPIclient, err := rockside.NewClientFromToken("token", "example.com", rockside.Testnet)
	if err != nil {
		panic(err)
	}

	smartWallets, err := rocksideAPIclient.SmartWallets.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(smartWallets)
}

func ExampleRPCClient() {
	// Get a RPC client from your existing Rockside client.
	rpc := rocksideClient.RPCClient

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

func ExampleClient_DeployContractWithSmartWallet() {
	smartWallets, err := rocksideClient.SmartWallets.List()
	if err != nil {
		panic(err)
	}

	var contractCode, jsonABI string
	txHash, err := rocksideClient.DeployContractWithSmartWallet(smartWallets[0], contractCode, jsonABI)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)
}

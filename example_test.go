package rockside_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocksideio/rockside-sdk-go"
)

var rocksideClient *rockside.Client

func ExampleNewClient() {
	rocksideAPIclient, err := rockside.NewClient(os.Getenv("ROCKSIDE_API_KEY"), rockside.Testnet)
	if err != nil {
		panic(err)
	}

	identities, err := rocksideAPIclient.Identities.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(identities)
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

func ExampleClient_DeployContractWithIdentity() {
	identities, err := rocksideClient.Identities.List()
	if err != nil {
		panic(err)
	}

	var contractCode, jsonABI string
	txHash, err := rocksideClient.DeployContractWithIdentity(identities[0], contractCode, jsonABI)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)
}

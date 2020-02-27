package rockside_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

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

func ExampleDeployContractWithIdentity() {
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

func Example_contractReadCallUsingContractABIBindings() {
	contractAddress := common.HexToAddress("my_contract_address")
	contract, err := NewContractCaller(contractAddress, rocksideClient.RPCClient)
	if err != nil {
		panic(err)
	}

	timestamp, _ := contract.Read(&bind.CallOpts{}, [32]byte{})
	fmt.Println(timestamp)
}

func Example_gaslessContractWriteCallUsingContractABIBindings() {
	rocksideIdentityAddress := common.HexToAddress("my_rockside_identity_contract_address")
	contractAddress := common.HexToAddress("my_contract_address")

	rocksideTransactor := rockside.NewTransactor(rocksideIdentityAddress, rocksideClient)
	contract, err := NewContractTransactor(contractAddress, rocksideTransactor)
	if err != nil {
		panic(err)
	}

	tx, _ := contract.Write(rockside.TransactOpts(), [32]byte{})

	txHash := rocksideTransactor.LookupRocksideTransactionHash(tx.Hash())
	fmt.Println(txHash)
}

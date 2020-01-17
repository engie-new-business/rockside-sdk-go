package rockside_test

import (
	"fmt"
	"github.com/rocksideio/rockside-sdk-go"
	"os"
)

var client *rockside.Client

func ExampleNewClient() {
	client, err := rockside.NewClient("https://api.rockside.io", os.Getenv("ROCKSIDE_API_KEY"))
	if err != nil {
		panic(err)
	}
	client.SetNetwork(rockside.Testnet)

	identities, err := client.Identities.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(identities)
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

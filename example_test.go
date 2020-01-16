package rockside_test

import (
	"fmt"
	rockside "github.com/rocksideio/rockside-sdk-go"
	"os"
)

func ExampleNewClient() {
	client, err := rockside.NewClient("https://api.rockside.io", os.Getenv("ROCKSIDE_API_KEY"))
	if err != nil {
		panic(err)
	}
	client.SetNetwork(rockside.Testnet)

	identities, _, err := client.Identities.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(identities)
}

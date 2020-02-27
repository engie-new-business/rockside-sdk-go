package rockside_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func Example_contractCaller() {
	contractAddress := common.HexToAddress("my_contract_address")

	// NewContractCaller is typically in your contract binding GO file which was generated via `abigen`
	contract, err := NewContractCaller(contractAddress, rocksideClient.RPCClient)
	if err != nil {
		panic(err)
	}

	timestamp, _ := contract.Read(&bind.CallOpts{}, [32]byte{})
	fmt.Println(timestamp)
}

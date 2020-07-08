package rockside

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	SmartWalletABI = `[{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"forwarder","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"value","type":"bytes"}],"name":"DataChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"salt","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"initCode","type":"bytes"}],"name":"Deployed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"destination","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"Executed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Received","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"success","type":"bool"}],"name":"RelayedExecute","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"},{"indexed":false,"internalType":"bool","name":"value","type":"bool"}],"name":"UpdateOwners","type":"event"},{"inputs":[],"name":"authorizedForwarder","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"struct SmartWallet.Call[]","name":"calls","type":"tuple[]"}],"name":"batch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes32","name":"salt","type":"bytes32"},{"internalType":"bytes","name":"initCode","type":"bytes"}],"name":"deploy","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"execute","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"}],"name":"getData","outputs":[{"internalType":"bytes","name":"value","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"forwarder","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"initialized","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"owners","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes","name":"value","type":"bytes"}],"name":"setData","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"bool","name":"value","type":"bool"}],"name":"updateOwners","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`
)

type SmartWallets endpoint

func (i *SmartWallets) Create(account, forwarder string) (ContractCreationResponse, error) {
	var result ContractCreationResponse

	path := fmt.Sprintf("ethereum/%s/smartwallets", i.client.network)
	req := struct {
		Account   string `json:"account"`
		Forwarder string `json:"forwarder"`
	}{account, forwarder}
	if _, err := i.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (i *SmartWallets) List() ([]string, error) {
	var result []string

	path := fmt.Sprintf("ethereum/%s/smartwallets", i.client.network)
	if _, err := i.client.get(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (i *SmartWallets) Exists(smartWalletAddr common.Address) (bool, error) {
	all, err := i.client.SmartWallets.List()
	if err != nil {
		return false, err
	}
	for _, item := range all {
		if item == smartWalletAddr.String() {
			return true, nil
		}
	}
	return false, nil
}

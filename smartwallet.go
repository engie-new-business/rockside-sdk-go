package rockside

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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

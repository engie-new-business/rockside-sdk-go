package rockside

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Identities endpoint

type createIdentitiesResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
}

func (i *Identities) Create() (createIdentitiesResponse, error) {
	var result createIdentitiesResponse

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	if _, err := i.client.post(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (i *Identities) List() ([]string, error) {
	var result []string

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	if _, err := i.client.get(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (i *Identities) Exists(identityAddress common.Address) (bool, error) {
	all, err := i.client.Identities.List()
	if err != nil {
		return false, err
	}
	for _, item := range all {
		if item == identityAddress.String() {
			return true, nil
		}
	}
	return false, nil
}

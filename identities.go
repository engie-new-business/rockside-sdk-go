package rockside

import (
	"fmt"
)

type identitiesEndpoint endpoint

type CreateIdentitiesResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
}

func (i *identitiesEndpoint) Create() (CreateIdentitiesResponse, error) {
	var result CreateIdentitiesResponse

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	if _, err := i.client.post(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (i *identitiesEndpoint) List() ([]string, error) {
	var result []string

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	if _, err := i.client.get(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

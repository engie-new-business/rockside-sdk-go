package rockside

import (
	"fmt"
	"net/http"
)

type IdentitiesEndpoint endpoint

type CreateIdentitiesResponse struct {
	Address string `json:"address"`
}

func (i *IdentitiesEndpoint) Create() (CreateIdentitiesResponse, *http.Response, error) {
	var result CreateIdentitiesResponse

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	resp, err := i.client.post(path, nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

func (i *IdentitiesEndpoint) List() ([]string, *http.Response, error) {
	var result []string

	path := fmt.Sprintf("ethereum/%s/identities", i.client.network)
	resp, err := i.client.get(path, nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

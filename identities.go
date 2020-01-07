package rockside

import "net/http"

type IdentitiesEndpoint endpoint

type CreateIdentitiesResponse struct {
	Address string `json:"address"`
}

func (i *IdentitiesEndpoint) Create(network Network) (CreateIdentitiesResponse, *http.Response, error) {

	result := CreateIdentitiesResponse{}

	resp, err := i.client.post("ethereum/"+network.String()+"/identities", nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

func (i *IdentitiesEndpoint) List(network Network) ([]string, *http.Response, error) {

	var result []string

	resp, err := i.client.get("ethereum/"+network.String()+"/identities", nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

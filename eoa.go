package rockside

import "net/http"

type EOAEndpoint endpoint

type CreateEOAResponse struct {
	Address string `json:"address"`
}

func (e *EOAEndpoint) Create() (CreateEOAResponse, *http.Response, error) {

	result := CreateEOAResponse{}

	resp, err := e.client.post("ethereum/eoa", nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

func (e *EOAEndpoint) List() ([]string, *http.Response, error) {

	var result []string

	resp, err := e.client.get("ethereum/eoa", nil, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

package rockside

type eoaEndpoint endpoint

type createEOAResponse struct {
	Address string `json:"address"`
}

func (e *eoaEndpoint) Create() (createEOAResponse, error) {
	var result createEOAResponse

	if _, err := e.client.post("ethereum/eoa", nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *eoaEndpoint) List() ([]string, error) {
	var result []string

	if _, err := e.client.get("ethereum/eoa", nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

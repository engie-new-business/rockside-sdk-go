package rockside

import (
	"fmt"
)

type Relay endpoint


func (e *Relay) GetParams(contractAddress string, account string) (RelayParamsResponse, error) {
	var result RelayParamsResponse
	path := fmt.Sprintf("ethereum/%s/relay/%s/params", e.client.network, contractAddress)
	_, err := e.client.get(path, nil, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}





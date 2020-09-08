package rockside

import (
	"fmt"
)

type Relay endpoint

type RelayTx struct {
	Data  string `json:"data"`
	Speed string `json:"speed"`
}

func (e *Relay) GetParams(destination string, account string) (RelayParamsResponse, error) {
	var result RelayParamsResponse
	path := fmt.Sprintf("ethereum/%s/relay/%s/params", e.client.network, destination)
	_, err := e.client.get(path, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (e *Relay) Relay(destination string, request RelayTx) (RelayTxResponse, error) {
	var result RelayTxResponse
	path := fmt.Sprintf("ethereum/%s/relay/%s", e.client.network, destination)
	if _, err := e.client.post(path, request, &result); err != nil {
		return result, err
	}
	return result, nil
}

package rockside

import (
	"fmt"
)

type RelayableIdentity endpoint

type RelayExecuteTxRequest struct {
	Relayer   string `json:"relayer"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	Gas       string `json:"gas"`
	GasPrice  string `json:"gasPrice"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type relayTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

type paramsResponse struct {
	Nonce   string `json:"nonce"`
	Relayer string `json:"relayer"`
}

type createRelayableIdentityResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

func (e *RelayableIdentity) Create(account string) (createRelayableIdentityResponse, error) {
	var result createRelayableIdentityResponse

	path := fmt.Sprintf("ethereum/%s/contracts/relayableidentity", e.client.network)
	req := struct {
		Account string `json:"account"`
	}{Account: account}
	if _, err := e.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

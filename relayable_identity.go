package rockside

import (
	"fmt"
)

type RelayableIdentity endpoint

type createRelayableIdentityResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

func (e *RelayableIdentity) Create(account, forwarder string) (createRelayableIdentityResponse, error) {
	var result createRelayableIdentityResponse

	path := fmt.Sprintf("ethereum/%s/contracts/relayableidentity", e.client.network)
	req := struct {
		Account   string `json:"account"`
		Forwarder string `json:"forwarder"`
	}{account, forwarder}
	if _, err := e.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

package rockside

import (
	"fmt"
)

type ContractsEndpoint endpoint

type createBouncerProxyRequest struct {
	Account string `json:"account"`
}

type CreateBouncerProxyResponse struct {
	BouncerProxyAddress string `json:"bouncer_proxy_address"`
	TransactionHash     string `json:"transaction_hash"`
}

func (c *ContractsEndpoint) CreateBouncerProxy(account string) (CreateBouncerProxyResponse, error) {
	var result CreateBouncerProxyResponse

	path := fmt.Sprintf("ethereum/%s/contracts/bouncerproxy", c.client.network)
	if _, err := c.client.post(path, createBouncerProxyRequest{Account: account}, &result); err != nil {
		return result, err
	}

	return result, nil
}

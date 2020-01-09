package rockside

import (
	"fmt"
	"net/http"
)

type ContractsEndpoint endpoint

type createBouncerProxyRequest struct {
	Account string `json:"account"`
}

type CreateBouncerProxyResponse struct {
	BouncerProxyAddress string `json:"bouncer_proxy_address"`
}

func (c *ContractsEndpoint) CreateBouncerProxy(account string) (CreateBouncerProxyResponse, *http.Response, error) {
	var result CreateBouncerProxyResponse

	path := fmt.Sprintf("ethereum/%s/contracts/bouncerproxy", c.client.network)
	resp, err := c.client.post(path, createBouncerProxyRequest{Account: account}, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

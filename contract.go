package rockside

import (
	"net/http"
)

type ContractsEndpoint endpoint

type createBouncerProxyRequest struct {
	Account string `json:"account"`
}

type CreateBouncerProxyResponse struct {
	BouncerProxyAddress string `json:"bouncer_proxy_address"`
}

func (c *ContractsEndpoint) CreateBouncerProxy(account string, network Network) (CreateBouncerProxyResponse, *http.Response, error) {
	result := CreateBouncerProxyResponse{}

	request := createBouncerProxyRequest{Account: account}

	resp, err := c.client.post("ethereum/"+network.String()+"/contracts/bouncerproxy", request, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil

}

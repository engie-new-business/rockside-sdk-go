package rockside

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewRPCClient(rocksideBaseURL, APIKey string, network Network) (*ethclient.Client, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/ethereum/%s/jsonrpc", rocksideBaseURL, network))
	if err != nil {
		return nil, fmt.Errorf("cannot build RPC URL from %s (%s)", rocksideBaseURL, network)
	}
	if endpoint.Scheme != "https" {
		return nil, fmt.Errorf("HTTPS scheme required for RPC client, got %s", rocksideBaseURL)
	}

	rpcClient, err := rpc.DialHTTPWithClient(endpoint.String(), &http.Client{Transport: &transport{APIKey}})
	if err != nil {
		return nil, fmt.Errorf("cannot RPC dial with custom HTTP client: %s", err)
	}

	return ethclient.NewClient(rpcClient), nil
}

type transport struct {
	apiKey string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("apikey", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

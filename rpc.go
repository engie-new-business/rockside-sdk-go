package rockside

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type RPCClient struct {
	endpoint       *url.URL
	authHTTPClient *http.Client
	*ethclient.Client
}

func newRPCClient(rocksideBaseURL, APIKey string, network Network) (*RPCClient, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/ethereum/%s/jsonrpc", rocksideBaseURL, network))
	if err != nil {
		return nil, fmt.Errorf("cannot build RPC URL from %s (%s)", rocksideBaseURL, network)
	}
	if endpoint.Scheme != "https" {
		return nil, fmt.Errorf("HTTPS scheme required for RPC client, got %s", rocksideBaseURL)
	}

	authHTTPClient := &http.Client{Transport: &transport{APIKey}}
	rpcClient, err := rpc.DialHTTPWithClient(endpoint.String(), authHTTPClient)
	if err != nil {
		return nil, fmt.Errorf("cannot RPC dial with custom HTTP client: %s", err)
	}

	return &RPCClient{
		endpoint:       endpoint,
		authHTTPClient: authHTTPClient,
		Client:         ethclient.NewClient(rpcClient),
	}, nil
}

const ethAccountsPayload = `{"jsonrpc":"2.0", "id": 1, "method": "eth_accounts", "params": []}`

func (r *RPCClient) EthAccounts() ([]string, error) {
	resp, err := r.authHTTPClient.Post(r.endpoint.String(), "", strings.NewReader(ethAccountsPayload))
	if err != nil {
		return []string{}, err
	}
	v := struct {
		Result []string `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return []string{}, err
	}
	return v.Result, nil
}

type transport struct {
	apiKey string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("apikey", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

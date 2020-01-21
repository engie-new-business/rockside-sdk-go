package rockside

import (
	"bytes"
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

type rpcRequest struct {
	ID      uint        `json:"id"`
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type rpcResponse struct {
	ID      uint        `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: %s (code=%d)", e.Message, e.Code)
}

func (r *RPCClient) SendTransactionFromIdentity(tx Transaction) (string, error) {
	if err := tx.ValidateFields(); err != nil {
		return "", err
	}

	accounts, err := r.EthAccounts()
	if err != nil {
		return "", nil
	}

	var found bool
	for _, a := range accounts {
		if strings.ToLower(a) == strings.ToLower(tx.From) {
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("transaction 'from' address '%s' is not one of your existing Rockside identities", tx.From)
	}

	body := &rpcRequest{ID: 1, Version: "2.0",
		Method: "eth_sendTransaction",
		Params: []Transaction{tx},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("cannot marshal RPC request to JSON: %s", err)
	}

	resp, err := r.authHTTPClient.Post(r.endpoint.String(), "", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return "", fmt.Errorf("cannot decode JSON RPC response: %s", err)
	}

	if e := rpcResp.Error; e != nil {
		return "", e
	}

	return fmt.Sprintf("%s", rpcResp.Result), nil
}

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

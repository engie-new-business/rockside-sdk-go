package rockside

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

type RPCClient struct {
	endpoint       *url.URL
	authHTTPClient *http.Client
	*ethclient.Client
}

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
	endpoint string
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: %s (code=%d, url=%s)", e.Message, e.Code, e.endpoint)
}

func (r *RPCClient) SendRocksideTransaction(tx Transaction) (string, error) {
	return r.sendTransaction(tx)
}

func (r *RPCClient) SendTransactionFromIdentity(tx Transaction) (string, error) {
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

	return r.sendTransaction(tx)
}

func (r *RPCClient) EthAccounts() ([]string, error) {
	body := &rpcRequest{ID: 1, Version: "2.0",
		Method: "eth_accounts",
		Params: []string{},
	}

	resp := struct {
		Result []string `json:"result"`
	}{}

	if err := r.post(body, &resp); err != nil {
		return []string{}, err
	}

	return resp.Result, nil
}

func (r *RPCClient) sendTransaction(tx Transaction) (string, error) {
	if err := tx.validateFields(); err != nil {
		return "", err
	}

	body := &rpcRequest{ID: 1, Version: "2.0",
		Method: "eth_sendTransaction",
		Params: []Transaction{tx},
	}

	resp := struct {
		Result string `json:"result"`
	}{}

	if err := r.post(body, &resp); err != nil {
		return "", err
	}

	return resp.Result, nil

}

func (r *RPCClient) post(body interface{}, decode interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("cannot marshal RPC request to JSON: %s", err)
	}

	resp, err := r.authHTTPClient.Post(r.endpoint.String(), "", bytes.NewReader(b))
	if err != nil {
		return err
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("cannot decode JSON RPC response: %s", err)
	}

	if err := rpcResp.Error; err != nil {
		err.endpoint = r.endpoint.String()
		return err
	}

	if decode != nil {
		b, err := json.Marshal(rpcResp)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, decode); err != nil {
			return err
		}
	}

	return nil
}

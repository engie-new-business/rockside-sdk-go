package rockside

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	Testnet Network = "ropsten"
	Mainnet Network = "mainnet"
)

type Network string

func (n Network) EtherscanURL() string {
	switch n {
	case Testnet:
		return fmt.Sprintf("https://%s.etherscan.io", n)
	default:
		return "https://etherscan.io"
	}
}

type endpoint struct {
	client *Client
}

type Client struct {
	baseURL *url.URL
	apiKey  string
	network Network

	client    *http.Client
	rpcClient *RPCClient
	logger    *log.Logger

	EOA          *EOAEndpoint
	Identities   *IdentitiesEndpoint
	Transaction  *TransactionEndpoint
	Contracts    *ContractsEndpoint
	BouncerProxy *BouncerProxyEndpoint
}

func NewClient(baseURL, APIKey string, net Network) (*Client, error) {
	var network Network

	switch net {
	case Mainnet, Testnet:
		network = net
	default:
		return nil, fmt.Errorf("init client: invalid network '%s' for client. Expecting: %s or %s", net, Mainnet, Testnet)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("init client: expected base URL with HTTPS scheme but got URL '%s'", baseURL)
	}

	if len(APIKey) == 0 {
		return nil, fmt.Errorf("init client: no API key found. Try with env variable ROCKSIDE_API_KEY")
	}
	if len(APIKey) != 32 {
		return nil, fmt.Errorf("init client: expected length of API key to be 32 but got %d", len(APIKey))
	}

	c := &Client{
		client:  http.DefaultClient,
		baseURL: u,
		apiKey:  APIKey,
		network: network,
		logger:  log.New(ioutil.Discard, "", 0),
	}

	c.EOA = &EOAEndpoint{c}
	c.Identities = &IdentitiesEndpoint{c}
	c.Transaction = &TransactionEndpoint{c}
	c.Contracts = &ContractsEndpoint{c}
	c.BouncerProxy = &BouncerProxyEndpoint{c}

	return c, nil
}

func (c *Client) SetLogger(l *log.Logger) {
	c.logger = l
}

func (c *Client) RPCClient() (*RPCClient, error) {
	if c.rpcClient == nil {
		return newRPCClient(c.baseURL.String(), c.apiKey, c.network)
	} else {
		return c.rpcClient, nil
	}
}

func (c *Client) CurrentNetwork() Network {
	return c.network
}

func (c *Client) URL() string {
	return c.baseURL.String()
}

func (c *Client) get(urlPath string, body interface{}, decode interface{}) (*http.Response, error) {
	return c.performRequest(http.MethodGet, urlPath, body, decode)
}

func (c *Client) post(urlPath string, body interface{}, decode interface{}) (*http.Response, error) {
	return c.performRequest(http.MethodPost, urlPath, body, decode)
}

func (c *Client) delete(urlPath string, body interface{}, decode interface{}) (*http.Response, error) {
	return c.performRequest(http.MethodDelete, urlPath, body, decode)
}

func (c *Client) put(urlPath string, body interface{}, decode interface{}) (*http.Response, error) {
	return c.performRequest(http.MethodPut, urlPath, body, decode)
}

func (c *Client) performRequest(method, urlPath string, body interface{}, decode interface{}) (*http.Response, error) {
	path, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	fullURL := c.baseURL.ResolveReference(path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, fullURL.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "rockside-go-client")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if tok := c.apiKey; len(tok) > 0 {
		req.Header.Set("apikey", tok)
	}

	dump, _ := httputil.DumpRequestOut(req, true)
	c.logger.Printf(">>>>>> Request %s-----\n\n", dump)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dump, _ = httputil.DumpResponse(resp, true)
	c.logger.Printf("<<<<<< Response %s\n-----\n\n", dump)

	if status := resp.StatusCode; status > 299 || status < 200 {
		if msg, err := decodeJSONErr(resp.Body); err != nil {
			c.logger.Printf("error body returned from '%s' does not seem to be JSON", resp.Request.URL)
			return resp, err
		} else {
			return resp, fmt.Errorf("%s (network: %s, URL: %s)", msg, c.CurrentNetwork(), c.URL())
		}
	}

	if decode != nil {
		if err := json.NewDecoder(resp.Body).Decode(decode); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func decodeJSONErr(body io.Reader) (string, error) {
	v := struct {
		Err     string `json:"error"`
		Message string `json:"message"`
	}{}
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		return "", err
	}
	if v.Err != "" {
		return v.Err, nil
	}
	return v.Message, nil
}

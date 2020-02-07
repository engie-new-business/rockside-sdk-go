package rockside

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

func (n Network) ChainID() *big.Int {
	switch n {
	case Testnet:
		return big.NewInt(3)
	default:
		return big.NewInt(1)
	}
}

type endpoint struct {
	client *Client
}

type Client struct {
	rocksideURL    *url.URL
	network        Network
	logger         *log.Logger
	authHTTPClient *http.Client

	RPCClient *RPCClient

	EOA               *EOAEndpoint
	Identities        *IdentitiesEndpoint
	Transaction       *TransactionEndpoint
	RelayableIdentity *RelayableIdentityEndpoint
}

func NewClient(rocksideBaseURL, apiKey string, net Network) (*Client, error) {
	var network Network

	switch net {
	case Mainnet, Testnet:
		network = net
	default:
		return nil, fmt.Errorf("init client: invalid network '%s' for client. Expecting: %s or %s", net, Mainnet, Testnet)
	}

	u, err := url.Parse(rocksideBaseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("init client: expected base URL with HTTPS scheme but got URL '%s'", rocksideBaseURL)
	}

	if len(apiKey) == 0 {
		return nil, fmt.Errorf("init client: no API key found. Try with env variable ROCKSIDE_API_KEY")
	}
	if len(apiKey) != 32 {
		return nil, fmt.Errorf("init client: expected length of API key to be 32 but got %d", len(apiKey))
	}

	rpcEndpoint, err := url.Parse(fmt.Sprintf("%s/ethereum/%s/jsonrpc", rocksideBaseURL, network))
	if err != nil {
		return nil, fmt.Errorf("cannot build RPC URL from %s (%s)", rocksideBaseURL, network)
	}

	authHTTPClient := &http.Client{Transport: &authTransport{apiKey}}
	rpcClient, err := rpc.DialHTTPWithClient(rpcEndpoint.String(), authHTTPClient)
	if err != nil {
		return nil, fmt.Errorf("cannot RPC dial with custom HTTP client: %s", err)
	}

	c := &Client{
		RPCClient: &RPCClient{
			endpoint:       rpcEndpoint,
			authHTTPClient: authHTTPClient,
			Client:         ethclient.NewClient(rpcClient),
		},
		authHTTPClient: authHTTPClient,
		rocksideURL:    u,
		network:        network,
		logger:         log.New(ioutil.Discard, "", 0),
	}

	c.EOA = &EOAEndpoint{c}
	c.Identities = &IdentitiesEndpoint{c}
	c.Transaction = &TransactionEndpoint{c}
	c.RelayableIdentity = &RelayableIdentityEndpoint{c}

	return c, nil
}

func (c *Client) SetLogger(l *log.Logger) {
	c.logger = l
}

func (c *Client) CurrentNetwork() Network {
	return c.network
}

func (c *Client) URL() string {
	return c.rocksideURL.String()
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
	fullURL := c.rocksideURL.ResolveReference(path)

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

	dump, _ := httputil.DumpRequestOut(req, true)
	c.logger.Printf(">>>>>> Request %s-----\n\n", dump)

	resp, err := c.authHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dump, _ = httputil.DumpResponse(resp, true)
	c.logger.Printf("<<<<<< Response %s\n-----\n\n", dump)

	if status := resp.StatusCode; status > 299 || status < 200 {
		context := fmt.Sprintf("[network: %s, URL: '%s']", c.CurrentNetwork(), c.URL())
		if msg, err := decodeJSONErr(resp.Body); err != nil {
			return resp, fmt.Errorf("body returned from %s does not seem to be JSON (try verbose mode): %s", context, err)
		} else {
			return resp, fmt.Errorf("%s %s", msg, context)
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

type authTransport struct {
	apiKey string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("apikey", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

package rockside

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	gethSigner "github.com/ethereum/go-ethereum/signer/core"
)

type Forwarder endpoint

type CreateForwarderResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

type RelayParamsResponse struct {
	Nonce     string            `json:"nonce"`
	GasPrices map[string]string `json:"gas_prices"`
}

type RelayExecuteTxData struct {
	Signer string `json:"signer"`
	To     string `json:"to"`
	Value  string `json:"value"`
	Data   string `json:"data"`
	Nonce  string `json:"nonce"`
}

type RelayExecuteTxRequest struct {
	DestinationContract string             `json:"destination_contract"`
	Speed               string             `json:"speed"`
	GasPriceLimit       string             `json:"gas_price_limit"`
	Data                RelayExecuteTxData `json:"data"`
	Signature           string             `json:"signature"`
}

type RelayTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

func (e *Forwarder) Create(owner string) (CreateForwarderResponse, error) {
	req := struct {
		Owner string `json:"owner"`
	}{owner}

	var result CreateForwarderResponse
	path := fmt.Sprintf("ethereum/%s/forwarders", e.client.network)
	if _, err := e.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *Forwarder) Get() ([]string, error) {
	var result []string
	path := fmt.Sprintf("ethereum/%s/forwarders", e.client.network)
	if _, err := e.client.get(path, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *Forwarder) GetRelayParams(forwarderAddress string, account string, channels ...string) (RelayParamsResponse, error) {
	channel := "0"
	if len(channels) > 0 {
		channel = channels[0]
	}
	var result RelayParamsResponse

	path := fmt.Sprintf("ethereum/%s/forwarders/%s/relayParams", e.client.network, forwarderAddress)
	req := struct {
		Account   string `json:"account"`
		ChannelID string `json:"channel_id"`
	}{Account: account, ChannelID: channel}
	_, err := e.client.post(path, req, &result)
	if err != nil {
		return result, err
	}

	channelNonce, isValidNonce := new(big.Int).SetString(result.Nonce, 10)
	if !isValidNonce {
		return RelayParamsResponse{}, fmt.Errorf("nonce is not a valid number [%s]", result.Nonce)
	}
	channelBig, isValidChannel := new(big.Int).SetString(channel, 10)
	if !isValidChannel {
		return RelayParamsResponse{}, fmt.Errorf("channel is not a valid number [%s]", channel)
	}
	return RelayParamsResponse{
		Nonce: new(big.Int).Add(new(big.Int).Lsh(channelBig, 128), channelNonce).String(),
	}, nil
}

func (e *Forwarder) Relay(forwarderAddress string, request RelayExecuteTxRequest) (RelayTxResponse, error) {
	var result RelayTxResponse

	if request.Speed == "" {
		request.Speed = "standard"
	}

	path := fmt.Sprintf("ethereum/%s/forwarders/%s", e.client.network, forwarderAddress)
	if _, err := e.client.post(path, request, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *Forwarder) SignTxParams(privateKeyStr, bouncerAddress, signer, destination, value, data, nonce string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	valueInt, ok := math.ParseBig256(value)
	if !ok {
		return "", errors.New("error with provided value")
	}

	if nonce == "" {
		paramsResponse, err := e.GetRelayParams(bouncerAddress, signer)
		if err != nil {
			return "", err
		}
		if nonce == "" {
			nonce = paramsResponse.Nonce
		}
	}
	nonceBig, isValidNonce := new(big.Int).SetString(nonce, 10)
	if !isValidNonce {
		return "", fmt.Errorf("nonce is not a valid number [%s]", nonce)
	}

	network := e.client.CurrentNetwork()

	argsHash, err := getHash(common.HexToAddress(bouncerAddress), common.HexToAddress(signer), common.HexToAddress(destination), valueInt, common.FromHex(data), nonceBig, network.ChainID())
	if err != nil {
		return "", err
	}

	signedHash, err := crypto.Sign(argsHash, privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signedHash), nil
}

func getHash(identity, signer, destination common.Address, value *big.Int, data []byte, nonce *big.Int, chainID *big.Int) ([]byte, error) {
	EIP712DomainType := []gethSigner.Type{
		{Name: "verifyingContract", Type: "address"},
		{Name: "chainId", Type: "uint256"},
	}

	txMessageType := []gethSigner.Type{
		{Name: "signer", Type: "address"},
		{Name: "to", Type: "address"},
		{Name: "value", Type: "uint256"},
		{Name: "data", Type: "bytes"},
		{Name: "nonce", Type: "uint256"},
	}

	types := gethSigner.Types{
		"TxMessage":    txMessageType,
		"EIP712Domain": EIP712DomainType,
	}

	domainData := gethSigner.TypedDataDomain{
		VerifyingContract: identity.String(),
		ChainId:           math.NewHexOrDecimal256(chainID.Int64()),
	}

	messageData := gethSigner.TypedDataMessage{
		"signer": signer.String(),
		"to":     destination.String(),
		"value":  value.String(),
		"data":   data,
		"nonce":  nonce.String(),
	}

	signerData := gethSigner.TypedData{
		Types:       types,
		PrimaryType: "TxMessage",
		Domain:      domainData,
		Message:     messageData,
	}

	typedDataHash, _ := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	domainSeparator, _ := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData), nil
}

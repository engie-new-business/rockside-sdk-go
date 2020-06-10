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

func (e *Forwarder) Create(owner string) (CreateForwarderResponse, error) {
	req := struct {
		Owner string `json:"owner"`
	}{owner}

	var result CreateForwarderResponse
	path := fmt.Sprintf("ethereum/%s/forwarder", e.client.network)
	if _, err := e.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *Forwarder) Get() (string, error) {
	type response struct {
		Address string `json:"address"`
	}

	var result response
	path := fmt.Sprintf("ethereum/%s/forwarder", e.client.network)
	if _, err := e.client.get(path, nil, &result); err != nil {
		return result.Address, err
	}

	return result.Address, nil
}

func (e *Forwarder) GetRelayParams(contractAddress string, account string, channels ...string) (paramsResponse, error) {
	channel := "0"
	if len(channels) > 0 {
		channel = channels[0]
	}
	var result paramsResponse

	path := fmt.Sprintf("ethereum/%s/%s/relayParams", e.client.network, contractAddress)
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
		return paramsResponse{}, fmt.Errorf("nonce is not a valid number [%s]", result.Nonce)
	}
	channelBig, isValidChannel := new(big.Int).SetString(channel, 10)
	if !isValidChannel {
		return paramsResponse{}, fmt.Errorf("channel is not a valid number [%s]", channel)
	}
	return paramsResponse{
		Nonce:   new(big.Int).Add(new(big.Int).Lsh(channelBig, 128), channelNonce).String(),
		Relayer: result.Relayer,
	}, nil
}

func (e *Forwarder) Relay(contractAddress string, request RelayExecuteTxRequest) (relayTxResponse, error) {
	var result relayTxResponse

	if request.Speed == "" {
		request.Speed = "standard"
	}

	path := fmt.Sprintf("ethereum/%s/%s/relay", e.client.network, contractAddress)
	if _, err := e.client.post(path, request, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *Forwarder) SignTxParams(privateKeyStr, bouncerAddress, relayer, signer, destination, value, data, gas, gasPrice, nonce string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	valueInt, ok := math.ParseBig256(value)
	if !ok {
		return "", errors.New("error with provided value")
	}

	gasUint, ok := math.ParseUint64(gas)
	if !ok {
		return "", errors.New("error with provided gas")
	}

	gasPriceInt, ok := math.ParseBig256(gasPrice)
	if !ok {
		return "", errors.New("error with provided gasPrice")
	}

	if nonce == "" || relayer == "" {
		paramsResponse, err := e.GetRelayParams(bouncerAddress, signer)
		if err != nil {
			return "", err
		}
		if nonce == "" {
			nonce = paramsResponse.Nonce
		}
		if relayer == "" {
			relayer = paramsResponse.Relayer
		}
	}
	nonceBig, isValidNonce := new(big.Int).SetString(nonce, 10)
	if !isValidNonce {
		return "", fmt.Errorf("nonce is not a valid number [%s]", nonce)
	}

	network := e.client.CurrentNetwork()

	argsHash, err := getHash(common.HexToAddress(bouncerAddress), common.HexToAddress(relayer), common.HexToAddress(signer), common.HexToAddress(destination), valueInt, common.FromHex(data), gasUint, gasPriceInt, nonceBig, network.ChainID())
	if err != nil {
		return "", err
	}

	signedHash, err := crypto.Sign(argsHash, privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signedHash), nil
}

func getHash(identity, relayer, signer, destination common.Address, value *big.Int, data []byte, gas uint64, gasPrice *big.Int, nonce *big.Int, chainID *big.Int) ([]byte, error) {
	EIP712DomainType := []gethSigner.Type{
		{Name: "verifyingContract", Type: "address"},
		{Name: "chainId", Type: "uint256"},
	}

	txMessageType := []gethSigner.Type{
		{Name: "relayer", Type: "address"},
		{Name: "signer", Type: "address"},
		{Name: "to", Type: "address"},
		{Name: "value", Type: "uint256"},
		{Name: "data", Type: "bytes"},
		{Name: "gasLimit", Type: "uint256"},
		{Name: "gasPrice", Type: "uint256"},
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
		"relayer":  relayer.String(),
		"signer":   signer.String(),
		"to":       destination.String(),
		"value":    value.String(),
		"data":     data,
		"gasLimit": fmt.Sprintf("%d", gas),
		"gasPrice": gasPrice.String(),
		"nonce":    nonce.String(),
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

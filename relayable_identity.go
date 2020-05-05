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

type RelayableIdentity endpoint

type RelayExecuteTxRequest struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	Gas       string `json:"gas"`
	GasPrice  string `json:"gasPrice"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type relayTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

type nonceResponse struct {
	Nonce string `json:"nonce"`
}

type createRelayableIdentityResponse struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

func (e *RelayableIdentity) Create(account string) (createRelayableIdentityResponse, error) {
	var result createRelayableIdentityResponse

	path := fmt.Sprintf("ethereum/%s/contracts/relayableidentity", e.client.network)
	req := struct {
		Account string `json:"account"`
	}{Account: account}
	if _, err := e.client.post(path, req, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *RelayableIdentity) RelayExecute(contractAddress string, request RelayExecuteTxRequest) (relayTxResponse, error) {
	var result relayTxResponse

	path := fmt.Sprintf("ethereum/%s/contracts/relayableidentity/%s/relayExecute", e.client.network, contractAddress)
	if _, err := e.client.post(path, request, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *RelayableIdentity) GetNonce(contractAddress string, account string, channels ...string) (nonceResponse, error) {
	channel := "0"
	if len(channels) > 0 {
		channel = channels[0]
	}
	var result nonceResponse

	path := fmt.Sprintf("ethereum/%s/contracts/relayableidentity/%s/nonce", e.client.network, contractAddress)
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
		return nonceResponse{}, fmt.Errorf("nonce is not a valid number [%s]", result.Nonce)
	}
	channelBig, isValidChannel := new(big.Int).SetString(channel, 10)
	if !isValidChannel {
		return nonceResponse{}, fmt.Errorf("channel is not a valid number [%s]", channel)
	}
	return nonceResponse{new(big.Int).Add(new(big.Int).Lsh(channelBig, 128), channelNonce).String()}, nil
}

func getHash(identity, signer, destination common.Address, value *big.Int, data []byte, gas uint64, gasPrice *big.Int, nonce *big.Int, chainID *big.Int) ([]byte, error) {
	EIP712DomainType := []gethSigner.Type{
		{Name: "verifyingContract", Type: "address"},
		{Name: "chainId", Type: "uint256"},
	}

	txMessageType := []gethSigner.Type{
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

func (b *RelayableIdentity) SignTxParams(privateKeyStr string, bouncerAddress string, signer string, destination string, value string, data string, gas string, gasPrice string, nonce string) (string, error) {
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

	gasPriceInt, ok := math.ParseBig256(value)
	if !ok {
		return "", errors.New("error with provided gasPrice")
	}

	if nonce == "" {
		nonceResponse, err := b.GetNonce(bouncerAddress, signer)
		if err != nil {
			return "", err
		}
		nonce = nonceResponse.Nonce
	}
	nonceBig, isValidNonce := new(big.Int).SetString(nonce, 10)
	if !isValidNonce {
		return "", fmt.Errorf("nonce is not a valid number [%s]", nonce)
	}

	network := b.client.CurrentNetwork()

	argsHash, err := getHash(common.HexToAddress(bouncerAddress), common.HexToAddress(signer), common.HexToAddress(destination), valueInt, common.FromHex(data), gasUint, gasPriceInt, nonceBig, network.ChainID())
	if err != nil {
		return "", err
	}

	signedHash, err := crypto.Sign(argsHash, privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signedHash), nil
}

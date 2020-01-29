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

type BouncerProxyEndpoint endpoint

type RelayTxRequest struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	Signature string `json:"signature"`
}

type RelayTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

type nonceRequest struct {
	Account string `json:"account"`
}

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

func (b *BouncerProxyEndpoint) Relay(contractAddress string, request RelayTxRequest) (RelayTxResponse, error) {
	var result RelayTxResponse

	path := fmt.Sprintf("ethereum/%s/contracts/bouncerproxy/%s/relay", b.client.network, contractAddress)
	if _, err := b.client.post(path, request, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (b *BouncerProxyEndpoint) GetNonce(contractAddress string, account string) (NonceResponse, error) {
	var result NonceResponse

	path := fmt.Sprintf("ethereum/%s/contracts/bouncerproxy/%s/nonce", b.client.network, contractAddress)
	_, err := b.client.post(path, nonceRequest{Account: account}, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func getHash(bouncer, signer, destination common.Address, value *big.Int, data []byte, nonce *big.Int, chainID *big.Int) ([]byte, error) {
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
		VerifyingContract: bouncer.String(),
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

func (b *BouncerProxyEndpoint) SignTxParams(privateKeyStr string, bouncerAddress string, signer string, destination string, value string, data string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	valueInt, ok := math.ParseBig256(value)
	if !ok {
		return "", errors.New("error with provided value")
	}

	nonceResponse, err := b.GetNonce(bouncerAddress, signer)
	if err != nil {
		return "", err
	}

	bouncerNonce := new(big.Int)
	bouncerNonce, _ = bouncerNonce.SetString(nonceResponse.Nonce, 10)

	network := b.client.CurrentNetwork()

	argsHash, err := getHash(common.HexToAddress(bouncerAddress), common.HexToAddress(signer), common.HexToAddress(destination), valueInt, common.FromHex(data), bouncerNonce, network.ChainID())
	if err != nil {
		return "", err
	}

	signedHash, err := crypto.Sign(argsHash, privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signedHash), nil
}

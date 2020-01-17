package rockside

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
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

func getHash(bouncer, signer, destination common.Address, value *big.Int, data []byte, nonce *big.Int) ([]byte, error) {
	addressTy, _ := abi.NewType("address", "", nil)
	uintTy, _ := abi.NewType("uint256", "", nil)
	bytesTy, _ := abi.NewType("bytes", "", nil)
	hashTy := abi.Arguments{
		abi.Argument{Name: "bouncer", Type: addressTy},
		abi.Argument{Name: "signer", Type: addressTy},
		abi.Argument{Name: "destination", Type: addressTy},
		abi.Argument{Name: "value", Type: uintTy},
		abi.Argument{Name: "data", Type: bytesTy},
		abi.Argument{Name: "nonce", Type: uintTy},
	}

	packed, err := hashTy.Pack(bouncer, signer, destination, value, data, nonce)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(packed), nil
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

	argsHash, err := getHash(common.HexToAddress(bouncerAddress), common.HexToAddress(signer), common.HexToAddress(destination), valueInt, common.FromHex(data), bouncerNonce)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(argsHash), argsHash)
	hash := crypto.Keccak256([]byte(msg))

	signedHash, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signedHash), nil
}

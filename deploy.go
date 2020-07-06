package rockside

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (c *Client) DeployContractWithSmartWallet(rocksideSmartWalletAddr, code, jsonABI string) (string, error) {
	if _, err := hexutil.Decode(rocksideSmartWalletAddr); err != nil {
		return "", fmt.Errorf("invalid smart wallet address: %s", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		return "", err
	}

	input, err := parsedABI.Pack("")
	if err != nil {
		return "", err
	}

	var data []byte
	data = append(common.FromHex(code), input...)
	resp, err := c.Transaction.Send(Transaction{
		From: rocksideSmartWalletAddr,
		Data: fmt.Sprintf("0x%x", data),
	})

	return resp.TransactionHash, err
}

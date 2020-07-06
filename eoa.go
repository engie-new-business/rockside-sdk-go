package rockside

import "fmt"

type EOA endpoint

func (e *EOA) Create() (addressResponse, error) {
	var result addressResponse

	if _, err := e.client.post("ethereum/eoa", nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (e *EOA) List() ([]string, error) {
	var result []string

	if _, err := e.client.get("ethereum/eoa", nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

type SignTransactionRequest struct {
	Transaction
	NetworkID string `json:"network_id"`
}

func (e *EOA) SignTransaction(address string, transaction SignTransactionRequest) (string, error) {
	path := fmt.Sprintf("ethereum/eoa/%s/sign", address)

	type signedTxResult struct {
		SignedTransaction string `json:"signed_transaction"`
	}

	var result signedTxResult
	if _, err := e.client.post(path, transaction, &result); err != nil {
		return "", err
	}

	return result.SignedTransaction, nil
}

type SignMessageRequest struct {
	Message string `json:"message"`
}

func (e *EOA) SignMessage(address string, message SignMessageRequest) (string, error) {
	path := fmt.Sprintf("ethereum/eoa/%s/sign-message", address)

	type signedTxResult struct {
		SignedMessage string `json:"signed_message"`
	}

	var result signedTxResult
	if _, err := e.client.post(path, message, &result); err != nil {
		return "", err
	}

	return result.SignedMessage, nil
}

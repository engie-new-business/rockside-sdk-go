package rockside

import (
	"fmt"
	"net/http"
)

type TransactionEndpoint endpoint

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Nonce    string `json:"nonce"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasprice"`
}

type SendTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

func (t *TransactionEndpoint) Send(transaction Transaction) (SendTxResponse, *http.Response, error) {
	var result SendTxResponse

	path := fmt.Sprintf("ethereum/%s/transaction", t.client.network)
	resp, err := t.client.post(path, transaction, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

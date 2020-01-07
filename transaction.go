package rockside

import "net/http"

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

func (t *TransactionEndpoint) Send(transaction Transaction, network Network) (SendTxResponse, *http.Response, error) {

	result := SendTxResponse{}

	resp, err := t.client.post("ethereum/"+network.String()+"/transaction", transaction, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

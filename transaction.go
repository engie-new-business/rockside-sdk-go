package rockside

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

	if err := validateTransactionFields(transaction); err != nil {
		return result, nil, err
	}

	path := fmt.Sprintf("ethereum/%s/transaction", t.client.network)
	resp, err := t.client.post(path, transaction, &result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

func validateTransactionFields(t Transaction) error {
	if err := validateHexField("from", t.From); err != nil {
		return err
	}
	if err := validateHexField("to", t.To); err != nil {
		return err
	}
	if err := validateHexField("data", t.Data); err != nil {
		return err
	}
	if err := validateHexField("value", t.Value); err != nil {
		return err
	}
	return nil
}

func validateHexField(fieldName, hexVal string) error {
	if len(hexVal) > 0 {
		if _, err := hexutil.Decode(hexVal); err != nil {
			return fmt.Errorf("invalid non empty '%s' field: %s", fieldName, err)
		}
	}
	return nil
}

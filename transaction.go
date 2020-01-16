package rockside

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
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
	if !common.IsHexAddress(t.From) {
		return errors.New("invalid 'from' address")
	}
	// To can be empty for contract creation
	if t.To != "" && !common.IsHexAddress(t.To) {
		return errors.New("invalid 'to' address")
	}
	if len(t.Data) > 0 {
		if _, err := hexutil.Decode(t.Data); err != nil {
			return fmt.Errorf("invalid 'data' bytes: %w", err)
		}
	}
	if len(t.Value) > 0 {
		if _, err := hexutil.DecodeBig(t.Value); err != nil {
			return fmt.Errorf("invalid 'value' number: %w", err)
		}
	}
	return nil
}

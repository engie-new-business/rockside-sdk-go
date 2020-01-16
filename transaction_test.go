package rockside

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestValidateTransactionFields(t *testing.T) {
	endpoint := new(TransactionEndpoint)
	tests := []struct {
		tx          Transaction
		errContains string
	}{
		{tx: Transaction{From: "123", To: ""}, errContains: "'from' field"},
		{tx: Transaction{From: "", To: "34567898"}, errContains: "'to' field"},
		{tx: Transaction{From: "1245", To: "34567898"}, errContains: "'from' field"},
		{tx: Transaction{Data: "456a789"}, errContains: "'data' field"},
		{tx: Transaction{Value: "456a789"}, errContains: "'value' field"},
	}

	for i, test := range tests {
		_, _, err := endpoint.Send(test.tx)
		if test.errContains == "" && err != nil {
			t.Fatalf("case %d: unexpected error %s", i+1, err)
		}
		if test.errContains != "" && err == nil {
			t.Fatalf("case %d: expected error, got none", i+1)
		}
		if sub := test.errContains; sub != "" && !strings.Contains(err.Error(), sub) {
			t.Fatalf("case %d: expecting error %q to contains %q", i+1, err, sub)
		}
	}
}

func TestTransaction(t *testing.T) {

	client, err := NewClient(baseURL, apikey)
	if err != nil {
		t.Fatal(err)
	}

	client.SetNetwork(Testnet)

	t.Run("Send transaction with Identity is OK", func(t *testing.T) {
		response, httpResponse, err := client.Identities.Create()
		if err != nil {
			t.Fatal(err)
		}

		blockDuration, err := strconv.Atoi(blockTime)
		if err != nil {
			t.Fatal(err)
		}

		//Need to wait for contract deployment's transaction to be mined
		time.Sleep(time.Duration(blockDuration) * time.Second)

		tx := Transaction{From: response.Address, To: response.Address, Value: "0x0"}
		txResponse, httpResponse, err := client.Transaction.Send(tx)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		if got, want := len(txResponse.TransactionHash), 66; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

package rockside

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestTransaction(t *testing.T) {

	client, err := New(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	client.SetAPIKey(apikey)

	t.Run("Send transaction with Identity is OK", func(t *testing.T) {
		response, httpResponse, err := client.Identities.Create(Mainnet)
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
		txResponse, httpResponse, err := client.Transaction.Send(tx, Mainnet)
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

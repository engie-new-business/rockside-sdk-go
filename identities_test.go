package rockside

import (
	"net/http"
	"testing"
)

func TestIdentities(t *testing.T) {

	client, err := NewClient(baseURL, apikey)
	if err != nil {
		t.Fatal(err)
	}

	client.SetNetwork(Testnet)

	t.Run("Create Identites", func(t *testing.T) {
		response, httpResponse, err := client.Identities.Create()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := httpResponse.StatusCode, http.StatusCreated; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		if got, want := len(response.Address), 42; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := len(response.TransactionHash), 66; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("List Identities", func(t *testing.T) {
		listResponse, httpResponse, err := client.Identities.List()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		initialNumberOfEOA := len(listResponse)
		if initialNumberOfEOA == 0 {
			t.Fatalf("expect response length %v greater than 0", initialNumberOfEOA)
		}

		createResponse, httpResponse, err := client.Identities.Create()
		if err != nil {
			t.Fatal(err)
		}

		listResponse, httpResponse, err = client.Identities.List()
		if err != nil {
			t.Fatal(err)
		}

		responseLenght := len(listResponse)
		if responseLenght <= initialNumberOfEOA {
			t.Fatalf("expect response length %v greater than %v", responseLenght, initialNumberOfEOA)
		}

		containsAddress := false
		for _, a := range listResponse {
			if a == createResponse.Address {
				containsAddress = true
			}
		}

		if !containsAddress {
			t.Fatalf("should contains created address")
		}
	})
}

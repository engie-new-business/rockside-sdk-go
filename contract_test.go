package rockside

import (
	"crypto/ecdsa"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	baseURL   = os.Getenv("ROCKSIDE_URL")
	apikey    = os.Getenv("ROCKSIDE_API_KEY")
	blockTime = os.Getenv("BLOCK_TIME")
)

func TestContract(t *testing.T) {
	client, err := NewClient(baseURL, apikey)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	t.Run("Create Bouncer proxy", func(t *testing.T) {

		bouncerResponse, httpResponse, err := client.Contracts.CreateBouncerProxy(fromAddress.String())
		if err != nil {
			t.Fatal(err)
		}

		if got, want := httpResponse.StatusCode, http.StatusCreated; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		if got, want := len(bouncerResponse.BouncerProxyAddress), 42; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		blockDuration, err := strconv.Atoi(blockTime)
		if err != nil {
			t.Fatal(err)
		}

		//Need to wait for contract deployment's transaction to be mined
		time.Sleep(time.Duration(blockDuration) * time.Second)

		t.Run("Bouncer proxy get nonce", func(t *testing.T) {

			response, httpResponse, err := client.BouncerProxy.GetNonce(bouncerResponse.BouncerProxyAddress, fromAddress.String())
			if err != nil {
				t.Fatal(err)
			}

			if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := response.Nonce, "0"; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}

		})

		t.Run("Bouncer proxy relay transaction", func(t *testing.T) {
			signature, err := client.BouncerProxy.SignTxParams(privateKeyString, bouncerResponse.BouncerProxyAddress, fromAddress.String(), fromAddress.String(), "0", "")
			if err != nil {
				t.Fatal(err)
			}

			request := RelayTxRequest{
				From:      fromAddress.String(),
				To:        fromAddress.String(),
				Signature: signature,
			}
			response, httpResponse, err := client.BouncerProxy.Relay(bouncerResponse.BouncerProxyAddress, request)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := httpResponse.StatusCode, http.StatusOK; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}

			if got, want := len(response.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	})

}

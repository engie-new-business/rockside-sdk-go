package rockside_test

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rocksideio/rockside-sdk-go"
)

/*
Launch the integration tests with: ROCKSIDE_API_URL=... ROCKSIDE_API_KEY=... go test -v
(or using BLOCK_WAIT_TIME env variable for a specific block wait time)
*/

var (
	blockWaitTime int
	rocksideURL   = os.Getenv("ROCKSIDE_API_URL")
)

func TestRockside(t *testing.T) {
	client, err := rockside.NewClientFromAPIKey(os.Getenv("ROCKSIDE_API_KEY"), rockside.Testnet, rocksideURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Identities", func(t *testing.T) {
		t.Parallel()

		t.Run("create", func(t *testing.T) {
			resp, err := client.Identities.Create()
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(resp.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := len(resp.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run("listing", func(t *testing.T) {
			listing, err := client.Identities.List()
			if err != nil {
				t.Fatal(err)
			}

			initialNumberOfEOA := len(listing)
			if initialNumberOfEOA == 0 {
				t.Fatalf("expect response length %v greater than 0", initialNumberOfEOA)
			}

			created, err := client.Identities.Create()
			if err != nil {
				t.Fatal(err)
			}

			listing, err = client.Identities.List()
			if err != nil {
				t.Fatal(err)
			}

			if l := len(listing); l <= initialNumberOfEOA {
				t.Fatalf("expect response length %v greater than %v", l, initialNumberOfEOA)
			}

			var hasAddr bool
			for _, a := range listing {
				if a == created.Address {
					hasAddr = true
				}
			}

			if !hasAddr {
				t.Fatalf("should contains created address")
			}
		})
	})

	t.Run("Transaction", func(t *testing.T) {
		t.Parallel()

		t.Run("Send transaction from identity", func(t *testing.T) {
			response, err := client.Identities.Create()
			if err != nil {
				t.Fatal(err)
			}

			//Need to wait for contract deployment's transaction to be mined
			time.Sleep(time.Duration(blockWaitTime) * time.Second)

			tx := rockside.Transaction{From: response.Address, To: response.Address, Value: "0x0"}
			txResponse, err := client.Transaction.Send(tx)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(txResponse.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

	})

	t.Run("EOA", func(t *testing.T) {
		t.Parallel()

		t.Run("create", func(t *testing.T) {
			resp, err := client.EOA.Create()
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(resp.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run("listing", func(t *testing.T) {
			listing, err := client.EOA.List()
			if err != nil {
				t.Fatal(err)
			}

			initialNumberOfEOA := len(listing)
			if initialNumberOfEOA == 0 {
				t.Fatalf("expect response length %v greater than 0", initialNumberOfEOA)
			}

			created, err := client.EOA.Create()
			if err != nil {
				t.Fatal(err)
			}

			listing, err = client.EOA.List()
			if err != nil {
				t.Fatal(err)
			}

			if l := len(listing); l <= initialNumberOfEOA {
				t.Fatalf("expect response length %v greater than %v", l, initialNumberOfEOA)
			}

			var hasAddr bool
			for _, a := range listing {
				if a == created.Address {
					hasAddr = true
				}
			}

			if !hasAddr {
				t.Fatalf("should contains created address")
			}
		})
	})

	t.Run("Contract", func(t *testing.T) {
		t.Parallel()

		privateKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyString := hexutil.Encode(privateKeyBytes)[2:]

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			t.Fatal("error casting public key to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		t.Run("Create relayable identity", func(t *testing.T) {
			relayableIdentity, err := client.RelayableIdentity.Create(fromAddress.String())
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(relayableIdentity.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := len(relayableIdentity.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}

			//Need to wait for contract deployment's transaction to be mined
			time.Sleep(time.Duration(blockWaitTime) * time.Second)

			t.Run("Relayable identity get nonce", func(t *testing.T) {
				resp, err := client.RelayableIdentity.GetRelayParams(relayableIdentity.Address, fromAddress.String())
				if err != nil {
					t.Fatal(err)
				}

				if got, want := resp.Nonce, "0"; got != want {
					t.Fatalf("got %v, want %v", got, want)
				}

				if resp.Relayer == "0" || resp.Relayer == "0x0000000000000000000000000000000000000000" {
					t.Fatalf("got empty relayer")
				}
			})

			t.Run("Relayable identity relay transaction", func(t *testing.T) {
				params, err := client.RelayableIdentity.GetRelayParams(relayableIdentity.Address, fromAddress.String())
				if err != nil {
					t.Fatal(err)
				}
				signature, err := client.RelayableIdentity.SignTxParams(privateKeyString, relayableIdentity.Address, params.Relayer, fromAddress.String(), fromAddress.String(), "0", "", "0", "0", params.Nonce)
				if err != nil {
					t.Fatal(err)
				}

				request := rockside.RelayExecuteTxRequest{
					Relayer:   params.Relayer,
					From:      fromAddress.String(),
					To:        fromAddress.String(),
					Signature: signature,
					Nonce:     params.Nonce,
				}
				resp, err := client.RelayableIdentity.RelayExecute(relayableIdentity.Address, request)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := len(resp.TransactionHash), 66; got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			})
		})
	})
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func init() {
	if len(rocksideURL) == 0 {
		exit("missing ROCKSIDE_API_URL env variable")
	}

	waitTime, exists := os.LookupEnv("BLOCK_WAIT_TIME")
	if !exists {
		waitTime = "120"
	}

	int, err := strconv.Atoi(waitTime)
	if err != nil {
		exit("cannot parse block wait time as int")
	}

	blockWaitTime = int

	fmt.Fprint(os.Stdout, fmt.Sprintf("Launching integration test on %s (block wait time %d)\n\n", rocksideURL, blockWaitTime))
}

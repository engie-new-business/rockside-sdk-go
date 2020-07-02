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
			eoa, err := client.EOA.Create()
			if err != nil {
				t.Fatal(err)
			}

			forwarder, err := client.Forwarder.Create(eoa.Address)
			if err != nil {
				t.Fatal(err)
			}

			identity, err := client.Identities.Create(eoa.Address, forwarder.Address)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(identity.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := len(identity.TransactionHash), 66; got != want {
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

			eoa, err := client.EOA.Create()
			if err != nil {
				t.Fatal(err)
			}

			forwarder, err := client.Forwarder.Create(eoa.Address)
			if err != nil {
				t.Fatal(err)
			}

			created, err := client.Identities.Create(eoa.Address, forwarder.Address)
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
			eoa, err := client.EOA.Create()
			if err != nil {
				t.Fatal(err)
			}

			forwarder, err := client.Forwarder.Create(eoa.Address)
			if err != nil {
				t.Fatal(err)
			}

			identity, err := client.Identities.Create(eoa.Address, forwarder.Address)
			if err != nil {
				t.Fatal(err)
			}

			//Need to wait for contract deployment's transaction to be mined
			time.Sleep(time.Duration(blockWaitTime) * time.Second)

			tx := rockside.Transaction{From: identity.Address, To: identity.Address, Value: "0x0"}
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

	t.Run("Forwarder contract", func(t *testing.T) {
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

		t.Run("deploy forwarder", func(t *testing.T) {
			forwarder, err := client.Forwarder.Create(fromAddress.String())
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(forwarder.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := len(forwarder.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}

			identity, err := client.Identities.Create(fromAddress.String(), forwarder.Address)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(identity.Address), 42; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			if got, want := len(identity.TransactionHash), 66; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}

			//Need to wait for contract deployment's transaction to be mined
			time.Sleep(time.Duration(blockWaitTime) * time.Second)

			t.Run("List relay params", func(t *testing.T) {
				resp, err := client.Forwarder.GetRelayParams(forwarder.Address, fromAddress.String())
				if err != nil {
					t.Fatal(err)
				}

				if got, want := resp.Nonce, "0"; got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			})

			t.Run("Forward transaction", func(t *testing.T) {
				params, err := client.Forwarder.GetRelayParams(forwarder.Address, fromAddress.String())
				if err != nil {
					t.Fatal(err)
				}
				signature, err := client.Forwarder.SignTxParams(privateKeyString, identity.Address, fromAddress.String(), "0x0000000000000000000000000000000000000000", "0", "", params.Nonce)
				if err != nil {
					t.Fatal(err)
				}

				request := rockside.RelayExecuteTxRequest{
					DestinationContract: identity.Address,
					Speed:               "standard",
					GasPriceLimit:       "30000000000",
					Signature:           signature,
					Data: rockside.RelayExecuteTxData{
						Signer: fromAddress.String(),
						To:     "0x0000000000000000000000000000000000000000",
						Value:  "0",
						Data:   "",
						Nonce:  params.Nonce,
					},
				}

				resp, err := client.Forwarder.Forward(forwarder.Address, request)
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

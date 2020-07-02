package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/rocksideio/rockside-sdk-go"
	"github.com/spf13/cobra"
)

var (
	forwarderCmd = &cobra.Command{
		Use:   "forwarder",
		Short: "Manage forwarders",
	}

	listForwardersCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List your forwarders",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := RocksideClient().Forwarder.List()
			if err != nil {
				return err
			}
			return printJSON(all)
		},
	}

	createForwarderCmd = &cobra.Command{
		Use:   "create",
		Short: "Deploy a new forwarder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing public address of the account")
			}

			forwarder, err := RocksideClient().Forwarder.Create(args[0])
			if err != nil {
				return err
			}

			return printJSON(forwarder)
		},
	}

	getForwarderParamsCmd = &cobra.Command{
		Use:   "params",
		Short: "Get forwarders params (nonce, gas prices, etc.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing contract address as first param and account address as second param")
			}
			contractAddress := args[0]
			accountAddress := args[1]

			nonce, err := RocksideClient().Forwarder.GetRelayParams(contractAddress, accountAddress)
			if err != nil {
				return err
			}

			return printJSON(nonce)
		},
	}

	sendTestTransactionCmd = &cobra.Command{
		Use:   "send-test-transaction",
		Short: "Build, sign and send transaction to be forwarded",
		RunE: func(cmd *cobra.Command, args []string) error {
			if forwaderAddressFlag == "" {
				return errors.New("missing forwarder address")
			}
			if rocksideEOAFlag == "" {
				return errors.New("missing Rockside custodian EOA (use for signature)")
			}
			if rocksideIdentityFlag == "" {
				return errors.New("missing Rockside identity (use as destination contract)")
			}

			sign, err := buildSigner(rocksideEOAFlag, "")
			if err != nil {
				return err
			}

			params, err := RocksideClient().Forwarder.GetRelayParams(forwaderAddressFlag, rocksideEOAFlag)
			if err != nil {
				return err
			}

			tx := &rockside.Transaction{
				To:    rocksideEOAFlag,
				Value: "0x00",
				Data:  "0x00",
				Nonce: params.Nonce,
			}

			signature, err := RocksideClient().Forwarder.SignTxParams(sign, forwaderAddressFlag, tx.From, tx.To, tx.Value, tx.Data, tx.Nonce)
			if err != nil {
				return err
			}
			fmt.Println("signature:")
			printJSON(signature)

			relayTx := &rockside.RelayExecuteTxRequest{
				DestinationContract: rocksideIdentityFlag,
				Speed:               "standard",
				GasPriceLimit:       params.GasPrices["standard"],
				Data: rockside.RelayExecuteTxData{
					Signer: rocksideEOAFlag,
					To:     tx.To,
					Value:  tx.Value,
					Data:   tx.Data,
					Nonce:  tx.Nonce,
				},
				Signature: signature,
			}

			response, err := RocksideClient().Forwarder.Forward(forwaderAddressFlag, *relayTx)
			if err != nil {
				return err
			}

			return printJSON(response)

		},
	}

	signTransactionCmd = &cobra.Command{
		Use:   "sign",
		Short: "sign transaction and parameters to be forwarded",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing contract address and transaction payload {\"from\":\"\",\"to\":\"\", \"value\":\"\", \"data\":\"\" }")
			}

			contractAddress := args[0]
			txJSON := args[1]
			tx := &rockside.Transaction{}
			if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
				return err
			}

			sign, err := buildSigner(rocksideEOAFlag, privateKeyFlag)
			if err != nil {
				return err
			}

			response, err := RocksideClient().Forwarder.SignTxParams(sign, contractAddress, tx.From, tx.To, tx.Value, tx.Data, tx.Nonce)
			if err != nil {
				return err
			}

			return printJSON(response)
		},
	}

	forwardTransactionCmd = &cobra.Command{
		Use:   "forward",
		Short: "Forward a transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing forwarder address and transaction payload {\"from\":\"\", \"to\":\"\", \"value\":\"\", \"speed\":\"\", \"gas_price_limit\":\"\", \"data\":\"\", \"signature\":\"\"}")
			}

			contractAddress := args[0]
			txJSON := args[1]
			relayTx := &rockside.RelayExecuteTxRequest{}
			if err := json.Unmarshal([]byte(txJSON), relayTx); err != nil {
				return err
			}

			relayResponse, err := RocksideClient().Forwarder.Forward(contractAddress, *relayTx)
			if err != nil {
				return err
			}

			return printJSON(relayResponse)
		},
	}
)

func buildSigner(custodianPrivateKey, localPrivateKey string) (rockside.SignFunc, error) {
	if len(localPrivateKey) > 0 {
		return func(m []byte) ([]byte, error) {
			privateKey, err := crypto.HexToECDSA(privateKeyFlag)
			if err != nil {
				return nil, err
			}
			signed, err := crypto.Sign(m, privateKey)
			if err != nil {
				return nil, err
			}
			return signed, nil
		}, nil
	} else if len(custodianPrivateKey) > 0 {
		return func(m []byte) ([]byte, error) {
			signed, err := RocksideClient().EOA.SignMessage(rocksideEOAFlag, rockside.SignMessageRequest{fmt.Sprintf("0x%x", m)})
			if err != nil {
				return nil, fmt.Errorf("cannot sign message with Rockside EOA %q: %s", rocksideEOAFlag, err)
			}
			return []byte(signed), nil
		}, nil
	}
	return nil, errors.New("no signers provided")
}

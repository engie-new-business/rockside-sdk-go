package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/rocksideio/rockside-sdk-go"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Rockside client",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return cmd.Usage()
	},
	SilenceUsage: true,
}

var (
	eoaCmd = &cobra.Command{
		Use:   "eoa",
		Short: "Manage EOA",
	}

	listEOACmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List EOA",
		RunE: func(cmd *cobra.Command, args []string) error {
			eoaList, err := RocksideClient().EOA.List()
			if err != nil {
				return err
			}
			printJSON(eoaList)
			return nil
		},
	}

	createEOACmd = &cobra.Command{
		Use:   "create",
		Short: "Create an EOA",
		RunE: func(cmd *cobra.Command, args []string) error {
			eoa, err := RocksideClient().EOA.Create()
			if err != nil {
				return err
			}
			printJSON(eoa)
			return nil
		},
	}
)

var (
	identitiesCmd = &cobra.Command{
		Use:   "identities",
		Short: "Manage identities",
	}

	listIdentitiesCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List identities",
		Example: "list identities ropsten",
		RunE: func(cmd *cobra.Command, args []string) error {
			identitiesList, err := RocksideClient().Identities.List()
			if err != nil {
				return err
			}
			printJSON(identitiesList)
			return nil
		},
	}

	createIdentitiesCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			identity, err := RocksideClient().Identities.Create()
			if err != nil {
				return err
			}
			printJSON(identity)
			return nil
		},
	}
)

var (
	relayableIdentityCmd = &cobra.Command{
		Use:   "relayableidentity",
		Short: "Manage relayable identities",
	}

	deployRelayableIdentityCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy a relayable identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("You need to provide the public address of the account")
			}

			relayableidentity, err := RocksideClient().RelayableIdentity.Create(args[0])
			if err != nil {
				return err
			}
			printJSON(relayableidentity)
			return nil
		},
	}

	getNonceCmd = &cobra.Command{
		Use:   "nonce",
		Short: "get nonce of a relayable identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Provide contract address as first param and account address as second param")
			}
			contractAddress := args[0]
			accountAddress := args[1]

			nonce, err := RocksideClient().RelayableIdentity.GetNonce(contractAddress, accountAddress)
			if err != nil {
				return err
			}
			printJSON(nonce)
			return nil
		},
	}

	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "sign transaction and parameters to be relayed",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Provide your contract address and tx {\"from\":\"\",\"to\":\"\", \"value\":\"\", \"data\":\"\" } as parameter.")
			}

			contractAddress := args[0]
			txJSON := args[1]
			tx := &rockside.Transaction{}
			if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
				return err
			}

			signResponse, err := RocksideClient().RelayableIdentity.SignTxParams(privateKeyFlag, contractAddress, tx.From, tx.To, tx.Value, tx.Data)

			if err != nil {
				return err
			}

			printJSON(signResponse)
			return nil
		},
	}

	relayCmd = &cobra.Command{
		Use:   "relay",
		Short: "relay transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Provide your contact address and tx parameters '{\"from\":\"\", \"to\":\"\" \"value\":\"\", \"data\":\"\", \"signature\":\"\"}'")
			}

			contractAddress := args[0]
			txJSON := args[1]
			relayTx := &rockside.RelayExecuteTxRequest{}
			if err := json.Unmarshal([]byte(txJSON), relayTx); err != nil {
				return err
			}

			relayResponse, err := RocksideClient().RelayableIdentity.RelayExecute(contractAddress, *relayTx)
			if err != nil {
				return err
			}

			printJSON(relayResponse)
			return nil
		},
	}
)

var (
	transactionCmd = &cobra.Command{
		Use:   "transaction",
		Short: "Manage transaction",
	}

	sentTxCmd = &cobra.Command{
		Use:   "send",
		Short: "send transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Provide your tx {\"from\":\"\",\"to\":\"\", \"value\":\"\", gas:\"\", \"gasPrice\":\"\", \"nonce\":\"\"} as parameter.")

			}
			txJSON := args[0]
			tx := &rockside.Transaction{}
			if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
				return err
			}

			txResponse, err := RocksideClient().Transaction.Send(*tx)
			if err != nil {
				return err
			}

			printJSON(txResponse)
			return nil
		},
	}
)

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	if err := enc.Encode(v); err != nil {
		log.Fatal(err)
	}
}

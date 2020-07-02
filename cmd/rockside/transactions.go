package main

import (
	"encoding/json"
	"errors"

	"github.com/rocksideio/rockside-sdk-go"
	"github.com/spf13/cobra"
)

var (
	transactionCmd = &cobra.Command{
		Use:     "transaction",
		Aliases: []string{"tx"},
		Short:   "Manage transaction",
	}

	sentTxCmd = &cobra.Command{
		Use:   "send",
		Short: "send transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing transaction payload {\"from\":\"\",\"to\":\"\", \"value\":\"\", gas:\"\", \"gasPrice\":\"\", \"nonce\":\"\"}")
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

			return printJSON(txResponse)
		},
	}

	showTxCmd = &cobra.Command{
		Use:   "show",
		Short: "show transaction given a tx hash or tracking ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing transaction hash or tracking ID")
			}
			result, err := RocksideClient().Transaction.Show(args[0])
			if err != nil {
				return err
			}

			return printJSON(result)
		},
	}

	listTxCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List your transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := RocksideClient().Transaction.List()
			if err != nil {
				return err
			}
			return printJSON(result)
		},
	}
)

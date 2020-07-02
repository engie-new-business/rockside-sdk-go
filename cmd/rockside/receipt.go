package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	showReceiptCmd = &cobra.Command{
		Use:   "receipt",
		Short: "List the transaction receipt for the given transaction hash",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing transaction hash")
			}

			txhash := common.HexToHash(args[0])

			receipt, err := RocksideClient().RPCClient.TransactionReceipt(context.Background(), txhash)
			if err != nil {
				return fmt.Errorf("with tx hash %s: %s", txhash.String(), err)
			}

			printJSON(receipt)
			return nil
		},
	}
)

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	rpcCmd = &cobra.Command{
		Use:   "rpc",
		Short: "Perform RPC call given as first argument",
		RunE: func(cmd *cobra.Command, args []string) error {
			rpc, err := RocksideClient().RPCClient()
			if err != nil {
				return err
			}

			if len(args) < 1 {
				return errors.New("missing RPC method to call")
			}

			switch args[0] {
			case "eth_gasPrice":
				price, err := rpc.SuggestGasPrice(context.Background())
				if err != nil {
					return err
				}
				fmt.Printf("gas price: %s", price)
			case "eth_accounts":
				accounts, err := rpc.EthAccounts()
				if err != nil {
					return err
				}
				fmt.Println(accounts)
			default:
				return fmt.Errorf("unknown RPC method %s", args[0])
			}
			return nil
		},
	}
)

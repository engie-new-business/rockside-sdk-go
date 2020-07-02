package main

import (
	"errors"

	"github.com/rocksideio/rockside-sdk-go"

	"github.com/spf13/cobra"
)

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

			return printJSON(eoaList)
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

			return printJSON(eoa)
		},
	}

	signMessageWithEOACmd = &cobra.Command{
		Use:   "sign-message",
		Short: "Sign your message using the given EAO",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing EOA public address as first param and message hexadecimal as second param")
			}
			account, message := args[0], args[1]
			signed, err := RocksideClient().EOA.SignMessage(account, rockside.SignMessageRequest{message})
			if err != nil {
				return err
			}

			return printJSON(signed)
		},
	}
)

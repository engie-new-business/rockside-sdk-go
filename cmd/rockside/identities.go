package main

import (
	"errors"

	"github.com/spf13/cobra"
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

			return printJSON(identitiesList)
		},
	}

	createIdentitiesCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a an identity given the account address and forwarder address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing public address of the account and/or forwarder address")
			}

			identity, err := RocksideClient().Identities.Create(args[0], args[1])
			if err != nil {
				return err
			}

			return printJSON(identity)
		},
	}
)

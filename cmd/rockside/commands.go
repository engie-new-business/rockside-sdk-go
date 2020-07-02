package main

import (
	"encoding/json"
	"errors"
	"os"

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
	tokensCmd = &cobra.Command{
		Use:   "tokens",
		Short: "Manage Tokens",
	}

	createTokenCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			domain := args[0]
			contracts := []string{}

			for i := 1; i < len(args); i++ {
				contracts = append(contracts, args[i])
			}

			token, err := RocksideClient().Tokens.Create(domain, contracts)
			if err != nil {
				return err
			}

			return printJSON(token)
		},
	}
)

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(v)
}

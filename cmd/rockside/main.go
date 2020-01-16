package main

import (
	"github.com/rocksideio/rockside-sdk-go"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	client            *rockside.Client
	envRocksideAPIKey = os.Getenv("ROCKSIDE_API_KEY")

	privateKeyFlag, rocksideURLFlag string
	testnetFlag, verboseFlag        bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&rocksideURLFlag, "url", "https://api.rockside.io", "Rockside API URL")
	rootCmd.PersistentFlags().BoolVar(&testnetFlag, "testnet", true, "Use testnet (Ropsten) instead of mainnet")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Verbose Rockside client")

	signCmd.PersistentFlags().StringVar(&privateKeyFlag, "privatekey", "", "privatekey")
	signCmd.MarkPersistentFlagRequired("privatekey")
	bouncerProxyCmd.AddCommand(deployBouncerProxyCmd, getNonceCmd, signCmd, relayCmd)
	transactionCmd.AddCommand(sentTxCmd)
	eoaCmd.AddCommand(listEOACmd, createEOACmd)
	identitiesCmd.AddCommand(listIdentitiesCmd, createIdentitiesCmd)

	rootCmd.AddCommand(eoaCmd, identitiesCmd, bouncerProxyCmd, transactionCmd)

	cobra.OnInitialize(func() {
		var err error
		client, err = rockside.NewClient(rocksideURLFlag, envRocksideAPIKey)
		if err != nil {
			log.Fatal(err)
		}
		if testnetFlag {
			client.SetNetwork(rockside.Testnet)
		}
		client.SetNetwork(rockside.Mainnet)

		if verboseFlag {
			client.SetLogger(log.New(os.Stderr, "", 0))
		}
	})
}

func main() {
	log.SetFlags(0)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

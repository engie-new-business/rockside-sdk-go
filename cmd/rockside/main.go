package main

import (
	"errors"
	"log"
	"os"

	"rockside/sdk-go"

	"github.com/spf13/cobra"
)

var (
	client      *rockside.Client
	rocksideURL string
	rocksideAPIKey  = os.Getenv("ROCKSIDE_API_KEY")
)

func init() {
	rootCmd.PersistentFlags().StringVar(&rocksideURL, "url", "https://api.rockside.io", "Rockside api URL")

	eoaCmd.AddCommand(listEOACmd, createEOACmd)

	listIdentitiesCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	listIdentitiesCmd.MarkPersistentFlagRequired("network")
	createIdentitiesCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	createIdentitiesCmd.MarkPersistentFlagRequired("network")
	identitiesCmd.AddCommand(listIdentitiesCmd, createIdentitiesCmd)

	deployBouncerProxyCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	deployBouncerProxyCmd.MarkPersistentFlagRequired("network")
	getNonceCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	getNonceCmd.MarkPersistentFlagRequired("network")
	signCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	signCmd.MarkPersistentFlagRequired("network")
	signCmd.PersistentFlags().StringVar(&privateKey, "privatekey", "", "privatekey")
	signCmd.MarkPersistentFlagRequired("privatekey")
	relayCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	relayCmd.MarkPersistentFlagRequired("network")
	bouncerProxyCmd.AddCommand(deployBouncerProxyCmd, getNonceCmd, signCmd, relayCmd)

	sentTxCmd.PersistentFlags().StringVar(&network, "network", "", "Network")
	sentTxCmd.MarkPersistentFlagRequired("network")
	transactionCmd.AddCommand(sentTxCmd)

	rootCmd.AddCommand(eoaCmd, identitiesCmd, bouncerProxyCmd, transactionCmd)

	cobra.OnInitialize(func() {
		var err error
		client, err = rockside.New(rocksideURL)
		if err != nil {
			log.Fatal(err)
		}

		if len(rocksideAPIKey) > 0 {
			client.SetAPIKey(rocksideAPIKey)
		} else {
			log.Fatal(errors.New("You need to provide an API Key"))
		}

	})
}

func main() {
	log.SetFlags(0)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

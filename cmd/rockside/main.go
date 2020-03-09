package main

import (
	"log"
	"os"

	"github.com/rocksideio/rockside-sdk-go"
)

var (
	envRocksideAPIKey      = os.Getenv("ROCKSIDE_API_KEY")
	envRocksideToken       = os.Getenv("ROCKSIDE_TOKEN")
	envRocksideTokenOrigin = os.Getenv("ROCKSIDE_TOKEN_ORIGIN")
	envRocksideAPIURL      = os.Getenv("ROCKSIDE_API_URL")

	rocksideTokenOrigin, rocksideURLFlag                  string
	privateKeyFlag, identityToDeployContractFlag          string
	testnetFlag, verboseFlag                              bool
	printContractABIFlag, printContractRuntimeBinFlag     bool
	compileContractOnlyFlag, printContractCreationBinFlag bool
)

func init() {
	if envRocksideAPIURL == "" {
		envRocksideAPIURL = "https://api.rockside.io"
	}

	rootCmd.PersistentFlags().StringVar(&rocksideURLFlag, "url", envRocksideAPIURL, "Rockside API URL")
	rootCmd.PersistentFlags().StringVar(&rocksideTokenOrigin, "token-origin", envRocksideTokenOrigin, "Origin associated with token")
	rootCmd.PersistentFlags().BoolVar(&testnetFlag, "testnet", false, "Use testnet (Ropsten) instead of mainnet")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Verbose Rockside client")

	signCmd.PersistentFlags().StringVar(&privateKeyFlag, "privatekey", "", "privatekey")
	signCmd.MarkPersistentFlagRequired("privatekey")
	relayableIdentityCmd.AddCommand(deployRelayableIdentityCmd, getNonceCmd, signCmd, relayCmd)
	transactionCmd.AddCommand(sentTxCmd)
	eoaCmd.AddCommand(listEOACmd, createEOACmd)
	identitiesCmd.AddCommand(listIdentitiesCmd, createIdentitiesCmd)
	tokensCmd.AddCommand(createTokenCmd)

	deployContractCmd.PersistentFlags().StringVar(&identityToDeployContractFlag, "identity-address", "", "Address of Rockside identity to use as 'from' when deploying contract")
	deployContractCmd.PersistentFlags().BoolVar(&printContractABIFlag, "print-abi", false, "Compile, print contract abi and exit")
	deployContractCmd.PersistentFlags().BoolVar(&printContractRuntimeBinFlag, "print-runtime-bin", false, "Compile, print contract runtime bytecode and exit")
	deployContractCmd.PersistentFlags().BoolVar(&printContractCreationBinFlag, "print-creation-bin", false, "Compile, print contract creation bytecode and exit")
	deployContractCmd.PersistentFlags().BoolVar(&compileContractOnlyFlag, "compile-only", false, "Compile without deploying and exit")

	rootCmd.AddCommand(eoaCmd, identitiesCmd, relayableIdentityCmd, transactionCmd, deployContractCmd, rpcCmd, showReceiptCmd, tokensCmd)
}

func RocksideClient() *rockside.Client {
	network := rockside.Mainnet
	if testnetFlag {
		network = rockside.Testnet
	}

	if envRocksideAPIKey != "" && envRocksideToken != "" {
		log.Fatal("both ROCKSIDE_API_KEY and ROCKSIDE_TOKEN are defined as environment variables. Pick one!")
	}

	var (
		client *rockside.Client
		err    error
	)
	if envRocksideAPIKey != "" {
		client, err = rockside.NewClientFromAPIKey(envRocksideAPIKey, network, rocksideURLFlag)
	}
	if envRocksideToken != "" {
		client, err = rockside.NewClientFromToken(envRocksideToken, rocksideTokenOrigin, network, rocksideURLFlag)
	}
	if err != nil {
		log.Fatal(err)
	}
	if verboseFlag {
		client.SetLogger(log.New(os.Stderr, "", 0))
	}
	return client
}

func main() {
	log.SetFlags(0)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

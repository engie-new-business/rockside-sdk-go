package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

var (
	deployContractCmd = &cobra.Command{
		Use:   "deploy-contract",
		Short: "Compile and deploy an Ethereum contract",
		Long:  "Given the filepath of a .sol contract (default .sol file in current dir), it will compile and deploy it using your Rockside smart wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			var solFilepath string

			if len(args) > 0 {
				solFilepath = args[0]
				if filepath.Ext(solFilepath) != ".sol" {
					return fmt.Errorf("expecting .sol file extension but got %q", solFilepath)
				}
			} else {
				var solidityFiles []string
				if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && filepath.Ext(info.Name()) == ".sol" {
						solidityFiles = append(solidityFiles, path)
					}
					return nil
				}); err != nil {
					return err
				}

				switch len(solidityFiles) {
				case 0:
					return errors.New("no contracts found in current directory")
				case 1:
					solFilepath = solidityFiles[0]
				default:
					return errors.New("multiple contracts in current directory. Give filepath of one contract")
				}
			}

			contracts, err := compileContracts(solFilepath)
			if err != nil {
				return err
			}

			if len(contracts) > 1 {
				return errors.New("error: multiple contract compiled result")
			}

			var contract *compiler.Contract
			for _, c := range contracts {
				contract = c
			}

			log.Printf("successfully compiled %q", solFilepath)

			if printContractABIFlag {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", " ")
				enc.Encode(contract.Info.AbiDefinition)
				fmt.Println()
			}

			if printContractCreationBinFlag {
				fmt.Println("\ncreation bytecode:", "\n", contract.Code)
				fmt.Println()
			}

			if printContractRuntimeBinFlag {
				fmt.Println("\nruntime bytecode:", "\n", contract.RuntimeCode)
				fmt.Println()
			}

			if compileContractOnlyFlag {
				os.Exit(0)
			}

			smartWallet := smartWalletToDeployContractFlag
			if smartWallet == "" {
				smartWallets, err := RocksideClient().SmartWallets.List()
				if err != nil {
					return fmt.Errorf("cannot list smart wallets: %s", err)
				}
				if len(smartWallets) == 0 {
					return errors.New("no Rockside smart wallet found")
				}

				smartWallet = smartWallets[len(smartWallets)-1]
			}

			b, err := json.Marshal(contract.Info.AbiDefinition)
			if err != nil {
				return fmt.Errorf("cannot marshal ABI JSON definition: %s", err)
			}

			log.Printf("deploying contract through Rockside smartWallet %s", smartWallet)

			tx, err := RocksideClient().DeployContractWithSmartWallet(smartWallet, contract.Code, string(b))
			if err != nil {
				return fmt.Errorf("cannot deploy contract: %s (txhash=%s)", err, tx)
			}

			log.Printf("successfully deployed contract with receipt %s/tx/%s", RocksideClient().CurrentNetwork().ExplorerURL(), tx)

			return nil
		},
	}
)

func compileContracts(file string) (map[string]*compiler.Contract, error) {
	return compiler.CompileSolidity("solc", file)
}

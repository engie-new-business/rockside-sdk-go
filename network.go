package rockside

import "errors"

type Network int

const (
	Mainnet Network = iota
	Ropsten
)

var (
	networksName = map[Network]string{Mainnet: "mainnet", Ropsten: "ropsten"}
)

func (d Network) String() string {
	return networksName[d]
}

func GetNetwork(network string) (Network, error) {
	for key, name := range networksName {
		if name == network {
			return key, nil
		}
	}
	return 0, errors.New("Invalid network")
}

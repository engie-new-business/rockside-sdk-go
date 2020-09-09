package rockside

type addressResponse struct {
	Address string `json:"address"`
}

type TransactionResponse struct {
	TransactionHash string `json:"transaction_hash"`
	TrackingID      string `json:"tracking_id"`
}

type ContractCreationResponse struct {
	addressResponse
	TransactionResponse
}

type paramsResponse struct {
	Nonce     string            `json:"nonce"`
	GasPrices map[string]string `json:"gas_prices"`
}

type RelayParamsResponse struct {
	Speeds map[string]SpeedInfo `json:"speeds"`
}

type SpeedInfo struct {
	GasPrice string `json:"gas_price"`
	Relayer  string `json:"relayer"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

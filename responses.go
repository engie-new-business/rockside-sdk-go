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
	speeds map[string] map[string]string `json:"speeds"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

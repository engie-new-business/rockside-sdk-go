package rockside

type Tokens endpoint

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (i *Tokens) Create(domain string, contracts []string) (CreateTokenResponse, error) {
	var result CreateTokenResponse
	req := struct {
		Domain    string   `json:"domain"`
		Contracts []string `json:"contracts"`
	}{Domain: domain, Contracts: contracts}

	if _, err := i.client.post("/tokens", req, &result); err != nil {
		return result, err
	}

	return result, nil
}

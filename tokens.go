package rockside

type Tokens endpoint

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (i *Tokens) Create(domain string, contracts []string) (CreateTokenResponse, error) {
	return i.createRequest(domain, "", contracts)
}

func (i *Tokens) CreateForEndUser(domain string, endUserID string, contracts []string) (CreateTokenResponse, error) {
	return i.createRequest(domain, endUserID, contracts)
}

func (i *Tokens) createRequest(origin string, endUserID string, contracts []string) (CreateTokenResponse, error) {
	var result CreateTokenResponse

	req := struct {
		Origin    string   `json:"origin"`
		EndUserID string   `json:"end_user_id"`
		Contracts []string `json:"contracts"`
	}{Origin: origin, EndUserID: endUserID, Contracts: contracts}

	if _, err := i.client.post("/tokens", req, &result); err != nil {
		return result, err
	}

	return result, nil
}

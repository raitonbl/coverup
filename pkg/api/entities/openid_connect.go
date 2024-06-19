package entities

type OpenIdClient struct {
	Secret   string   `json:"secret"`
	Scopes   []string `json:"scopes"`
	ClientId string   `json:"client_id"`
}

type OpenIdUsernameAndPassword struct {
	OpenIdClient
	User UsernameAndPassword `json:"user"`
}

package entities

type OpenIdClient struct {
	BasicEntity
	Secret   string   `json:"secret"`
	Scopes   []string `json:"scopes"`
	ClientId string   `json:"client_id"`
}

type OpenIdUsernameAndPassword struct {
	BasicEntity
	OpenIdClient
	User UsernameAndPassword `json:"user"`
}

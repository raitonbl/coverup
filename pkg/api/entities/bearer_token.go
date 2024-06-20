package entities

type BearerToken struct {
	BasicEntity
	Value string `json:"value"`
}

package api

type BearerToken struct {
	BasicEntity
	Value string `json:"value"`
}

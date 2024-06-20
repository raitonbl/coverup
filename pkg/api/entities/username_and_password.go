package entities

type UsernameAndPassword struct {
	BasicEntity
	Username string `json:"username"`
	Password string `json:"password"`
}

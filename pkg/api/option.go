package api

type Option struct {
	Regexp         string
	Description    string
	HandlerFactory HandlerFactory
}

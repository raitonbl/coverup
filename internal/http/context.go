package internal

type Context interface {
	GetServerURL() string
	GetWorkDirectory() string
	GetHttpClient() HttpClient
}

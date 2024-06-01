package internal

import "github.com/raitonbl/coverup/pkg"

type Context interface {
	GetServerURL() string
	GetWorkDirectory() string
	GetHttpClient() pkg.HttpClient
	GetResourcesHttpClient() pkg.HttpClient
}

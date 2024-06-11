package context

import "github.com/raitonbl/coverup/pkg"

type Context interface {
	GetServerURL() string
	GetHttpClient() pkg.HttpClient
	GetResourcesHttpClient() pkg.HttpClient
}

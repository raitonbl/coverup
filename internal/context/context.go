package context

import (
	"github.com/raitonbl/coverup/pkg/http"
)

type Context interface {
	GetServerURL() string
	GetHttpClient() http.Client
	GetResourcesHttpClient() http.Client
}

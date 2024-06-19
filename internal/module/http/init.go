package http

import (
	"github.com/raitonbl/coverup/pkg/api"
)

type URIScheme string

const (
	ComponentType = "HttpRequest"
)

const (
	noneUriScheme  URIScheme = ""
	fileUriScheme  URIScheme = "file"
	httpUriScheme  URIScheme = "http"
	httpsUriScheme URIScheme = "https"
)

const (
	serverURLRegex   = `(https?://[^\s]+)`
	relativeURIRegex = `/([^/]+(?:/[^/]+)*)`
	valueRegex       = api.ValueExpression
	httpRequestRegex = `\{\{\s*` + ComponentType + `\.[a-zA-Z0-9_]+\s*\}\}`
	propertyRegex    = api.PropertyExpression
)

func OnV3(ctx api.StepDefinitionContext) {
	arr := []api.StepFactory{
		&GivenHttpRequestStepFactory{},
		&ThenHttpResponseStepFactory{},
	}
	for _, each := range arr {
		each.New(ctx)
	}
}

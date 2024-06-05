package v3

import (
	"fmt"
	v3 "github.com/raitonbl/coverup/internal/v3"
)

const valueRegex = `\{\{\s*([a-zA-Z0-9]+\.)*[a-zA-Z0-9]+\s*\}\}`
const relativeURIRegex = `^(/|(\./|\.\./|[a-zA-Z0-9_\-\.]+/?)*[a-zA-Z0-9_\-\.]+/?)+`
const serverURLRegex = `^(http|https)://([a-zA-Z0-9_\-\.]+(:\d+)?(/.*)?)`

func Apply(ctx v3.ScenarioContext) {
	h := HttpContext{
		ctx: ctx,
	}
	// Define Request
	ctx.GerkhinContext().Given(`^(?i)a HttpRequest$`, h.WithRequest)
	ctx.GerkhinContext().Given(`(?i)^a HttpRequest named (.+)$`, h.WithRequestWhenAlias)
	// Define method
	ctx.GerkhinContext().Step(`^(?i)method is ([^"]*)$`, h.WithMethod)
	ctx.GerkhinContext().Step(`^(?i)the method is ([^"]*)$`, h.WithMethod)
	// Define headers
	ctx.GerkhinContext().Step(`^(?i)headers:$`, h.WithHeaders)
	ctx.GerkhinContext().Step(`^(?i)the headers:$`, h.WithHeaders)
	ctx.GerkhinContext().Step(`^(?i)header (.*) is "([^"]*)"$`, h.WithHeader)
	ctx.GerkhinContext().Step(`^(?i)the header (.*) is "([^"]*)"$`, h.WithHeader)
	ctx.GerkhinContext().Step(`^(?i)the Header (.*) is "([^"]*)"$`, h.WithHeader)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)header (.*) is "%s"$`, valueRegex), h.WithHeader)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the header (.*) is "%s"$`, valueRegex), h.WithHeader)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Header (.*) is "%s"$`, valueRegex), h.WithHeader)
	ctx.GerkhinContext().Step(`^(?i)accept is "([^"]*)"$`, h.WithAcceptHeader)
	ctx.GerkhinContext().Step(`^(?i)the accept is "([^"]*)"$`, h.WithAcceptHeader)
	ctx.GerkhinContext().Step(`^(?i)the Accept is "([^"]*)"$`, h.WithAcceptHeader)
	ctx.GerkhinContext().Step(`^(?i)content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	ctx.GerkhinContext().Step(`^(?i)the content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	ctx.GerkhinContext().Step(`^(?i)the Content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	// Server URL & Path
	ctx.GerkhinContext().Step(`^(?i)path is http://(.+)$`, h.WithHttpPath)
	ctx.GerkhinContext().Step(`^(?i)the path is http://(.+)$`, h.WithHttpPath)
	ctx.GerkhinContext().Step(`^(?i)the Path is http://(.+)$`, h.WithHttpPath)
	ctx.GerkhinContext().Step(`^(?i)path is https://(.+)$`, h.WithHttpsPath)
	ctx.GerkhinContext().Step(`^(?i)the path is https://(.+)$`, h.WithHttpsPath)
	ctx.GerkhinContext().Step(`^(?i)the Path is https://(.+)$`, h.WithHttpsPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)path is %s$`, valueRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the path is %s$`, valueRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Path is %s$`, valueRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)path is %s$`, relativeURIRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the path is %s$`, relativeURIRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Path is %s$`, relativeURIRegex), h.WithPath)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)server url is %s$`, serverURLRegex), h.WithServerURL)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the server url is %s$`, serverURLRegex), h.WithServerURL)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Server url is %s$`, serverURLRegex), h.WithServerURL)
	//Body
	ctx.GerkhinContext().Step(`^(?i)(the |)body is $`, h.WithBody)
	ctx.GerkhinContext().Step(`^(?i)(the |)body is file://(.+)$`, h.WithBodyFileURI)
	//Form
	ctx.GerkhinContext().Step(`^(?i)(the |)form enctype is ([^"]*)$`, h.WithFormEncType)
	ctx.GerkhinContext().Step(`^(?i)(the |)form attribute "([a-zA-Z_]+)" is "([^"]+)"$`, h.WithFormAttribute)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)(the |)form attribute "%s"$`, valueRegex), h.WithFormAttribute)
	ctx.GerkhinContext().Step(`^form attribute "([a-zA-Z_]+)" is file://(.+)$`, h.WithFormFile)
	ctx.GerkhinContext().Step(fmt.Sprintf(`^form attribute "([a-zA-Z_]+)" is file://%s`, valueRegex), h.WithFormFile)

}

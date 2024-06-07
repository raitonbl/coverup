package v3

import (
	"fmt"
	v3 "github.com/raitonbl/coverup/internal/v3"
)

const serverURLRegex = `(https?://[^\s]+)`
const relativeURIRegex = `/([^/]+(?:/[^/]+)*)`
const valueRegex = `\{\{\s*([a-zA-Z0-9]+\.)*[a-zA-Z0-9]+\s*\}\}`
const httpRequestRegex = `\{\{\s*HttpRequest\.\w+\s*\}\}`
const entityRegex = `\{\{\s*Entity\.\w+\s*\}\}`

func Apply(ctx v3.ScenarioContext) {
	h := &HttpContext{
		ctx: ctx,
	}
	onRequest(h)
	onResponse(h)
}

func onRequest(h *HttpContext) {
	// Define Request
	h.ctx.GerkhinContext().Given(`^(?i)a HttpRequest$`, h.WithRequest)
	h.ctx.GerkhinContext().Given(`(?i)^a HttpRequest named (.+)$`, h.WithRequestWhenAlias)
	// Define method
	h.ctx.GerkhinContext().Step(`^(?i)method is ([^"]*)$`, h.WithMethod)
	h.ctx.GerkhinContext().Step(`^(?i)the method is ([^"]*)$`, h.WithMethod)
	// Define headers
	h.ctx.GerkhinContext().Step(`^(?i)headers:$`, h.WithHeaders)
	h.ctx.GerkhinContext().Step(`^(?i)the headers:$`, h.WithHeaders)
	h.ctx.GerkhinContext().Step(`^(?i)header (.*) is "([^"]*)"$`, h.WithHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the header (.*) is "([^"]*)"$`, h.WithHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the Header (.*) is "([^"]*)"$`, h.WithHeader)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)header (.*) is "%s"$`, valueRegex), h.WithHeader)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the header (.*) is "%s"$`, valueRegex), h.WithHeader)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Header (.*) is "%s"$`, valueRegex), h.WithHeader)
	h.ctx.GerkhinContext().Step(`^(?i)accept is "([^"]*)"$`, h.WithAcceptHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the accept is "([^"]*)"$`, h.WithAcceptHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the Accept is "([^"]*)"$`, h.WithAcceptHeader)
	h.ctx.GerkhinContext().Step(`^(?i)content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	h.ctx.GerkhinContext().Step(`^(?i)the Content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	// Server URL & Path
	h.ctx.GerkhinContext().Step(`^(?i)path is http://(.+)$`, h.WithHttpPath)
	h.ctx.GerkhinContext().Step(`^(?i)the path is http://(.+)$`, h.WithHttpPath)
	h.ctx.GerkhinContext().Step(`^(?i)the Path is http://(.+)$`, h.WithHttpPath)
	h.ctx.GerkhinContext().Step(`^(?i)path is https://(.+)$`, h.WithHttpsPath)
	h.ctx.GerkhinContext().Step(`^(?i)the path is https://(.+)$`, h.WithHttpsPath)
	h.ctx.GerkhinContext().Step(`^(?i)the Path is https://(.+)$`, h.WithHttpsPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)path is %s$`, valueRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the path is %s$`, valueRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Path is %s$`, valueRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)path is %s$`, relativeURIRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the path is %s$`, relativeURIRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Path is %s$`, relativeURIRegex), h.WithPath)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)server url is %s$`, serverURLRegex), h.WithServerURL)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the server url is %s$`, serverURLRegex), h.WithServerURL)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Server url is %s$`, serverURLRegex), h.WithServerURL)
	//Body
	h.ctx.GerkhinContext().Step(`^(?i)body is:$`, h.WithBody)
	h.ctx.GerkhinContext().Step(`^(?i)the body is:$`, h.WithBody)
	h.ctx.GerkhinContext().Step(`^(?i)the Body is:$`, h.WithBody)
	h.ctx.GerkhinContext().Step(`^(?i)body is file://(.+)$`, h.WithBodyFileURI)
	h.ctx.GerkhinContext().Step(`^(?i)the body is file://(.+)$`, h.WithBodyFileURI)
	h.ctx.GerkhinContext().Step(`^(?i)the Body is file://(.+)$`, h.WithBodyFileURI)
	//Form
	h.ctx.GerkhinContext().Step(`^(?i)form enctype is ([^"]*)$`, h.WithFormEncType)
	h.ctx.GerkhinContext().Step(`^(?i)the form enctype is ([^"]*)$`, h.WithFormEncType)
	h.ctx.GerkhinContext().Step(`^(?i)the Form enctype is ([^"]*)$`, h.WithFormEncType)
	h.ctx.GerkhinContext().Step(`^(?i)form attribute "([a-zA-Z_]+)" is "([^"]+)"$`, h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(`^(?i)the form attribute "([a-zA-Z_]+)" is "([^"]+)"$`, h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(`^(?i)the Form attribute "([a-zA-Z_]+)" is "([^"]+)"$`, h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)form attribute "%s"$`, valueRegex), h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the form attribute "%s"$`, valueRegex), h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Form attribute "%s"$`, valueRegex), h.WithFormAttribute)
	h.ctx.GerkhinContext().Step(`^(?i)form attribute "([a-zA-Z_]+)" is file://(.+)$`, h.WithFormFile)
	h.ctx.GerkhinContext().Step(`^(?i)the form attribute "([a-zA-Z_]+)" is file://(.+)$`, h.WithFormFile)
	h.ctx.GerkhinContext().Step(`^(?i)the Form attribute "([a-zA-Z_]+)" is file://(.+)$`, h.WithFormFile)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)form attribute "([a-zA-Z_]+)" is file://%s`, valueRegex), h.WithFormFile)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the form attribute "([a-zA-Z_]+)" is file://%s`, valueRegex), h.WithFormFile)
	h.ctx.GerkhinContext().Step(fmt.Sprintf(`^(?i)the Form attribute "([a-zA-Z_]+)" is file://%s`, valueRegex), h.WithFormFile)
	// Submit
	h.ctx.GerkhinContext().When("^(?i)client submits the HttpRequest", h.SubmitHttpRequest)
	h.ctx.GerkhinContext().When("^(?i)the client submits the HttpRequest", h.SubmitHttpRequest)
	h.ctx.GerkhinContext().When("^(?i)the Client submits the HttpRequest", h.SubmitHttpRequest)
	h.ctx.GerkhinContext().When(fmt.Sprintf("%s submits the HttpRequest", entityRegex), h.SubmitHttpRequestOnBehalfOfEntity)
	h.ctx.GerkhinContext().When(fmt.Sprintf("^(?i)the %s submits the HttpRequest", entityRegex), h.SubmitHttpRequestOnBehalfOfEntity)
	h.ctx.GerkhinContext().When(fmt.Sprintf("^(?i)client submits %s", httpRequestRegex), h.SubmitNamedHttpRequest)
	h.ctx.GerkhinContext().When(fmt.Sprintf("^(?i)the client submits %s", httpRequestRegex), h.SubmitNamedHttpRequest)
	h.ctx.GerkhinContext().When(fmt.Sprintf("^(?i)the Client submits %s", httpRequestRegex), h.SubmitNamedHttpRequest)
	h.ctx.GerkhinContext().When(fmt.Sprintf("%s submits %s", entityRegex, httpRequestRegex), h.SubmitNamedHttpRequestOnBehalfOfEntity)
	h.ctx.GerkhinContext().When(fmt.Sprintf("^(?i)the %s submits %s", entityRegex, httpRequestRegex), h.SubmitNamedHttpRequestOnBehalfOfEntity)

}

func onResponse(h *HttpContext) {
	// Status code
	h.ctx.GerkhinContext().Then(`^(?i)response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(`^(?i)the response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(`^(?i)the Response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response status code for %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response status code %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response status code %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
	// headers
	h.ctx.GerkhinContext().Then(`^(?i)response headers are:$`, h.AssertResponseExactHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)the response headers are:$`, h.AssertResponseExactHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)the Response headers are:$`, h.AssertResponseExactHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)response headers contains:$`, h.AssertResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)the response headers contains:$`, h.AssertResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)the Response headers contains:$`, h.AssertResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(`^(?i)response content-type is "([^"]*)"$`, h.AssertResponseContentType)
	h.ctx.GerkhinContext().Then(`^(?i)the response content-type is "([^"]*)"$`, h.AssertResponseContentType)
	h.ctx.GerkhinContext().Then(`^(?i)the Response content-type is "(.+)"$`, h.AssertResponseContentType)
	h.ctx.GerkhinContext().Then(`^(?i)response header ([^ ]+) is "([^"]*)"$`, h.AssertResponseHeader)
	h.ctx.GerkhinContext().Then(`^(?i)the response header ([^ ]+) is "([^"]*)"$`, h.AssertResponseHeader)
	h.ctx.GerkhinContext().Then(`^(?i)the Response header is "([^"]*)"$`, h.AssertResponseHeader)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response headers for %s are:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseExactHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response headers for %s are:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseExactHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response headers for %s are:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseExactHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response headers for %s contains:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response headers for %s contains:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response headers for %s contains:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContainsHeaders)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response content-type for %s is "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContentType)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response content-type for %s is "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContentType)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response content-type for %s is "(.+)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseContentType)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response header ([^ ]+) for %s is "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseHeader)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response header ([^ ]+) for %s is "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseHeader)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response header for %s is "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponseHeader)
	// body:schema
	h.ctx.GerkhinContext().Then(`^(?i)response body respects schema file://(.+)$`, h.AssertResponseIsValidAgainstSchema)
	h.ctx.GerkhinContext().Then(`^(?i)the response body respects schema file://(.+)$`, h.AssertResponseIsValidAgainstSchema)
	h.ctx.GerkhinContext().Then(`^(?i)the Response body respects schema file://(.+)$`, h.AssertResponseIsValidAgainstSchema)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body for %s respects schema file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseIsValidAgainstSchema)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body for %s respects schema file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseIsValidAgainstSchema)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body for %s respects schema file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseIsValidAgainstSchema)
	// body
	patterns := map[string]any{
		`"([^"]*)"`:    h.AssertResponseBodyPathEqualsTo,
		valueRegex:     h.AssertResponseBodyPathEqualsToValue,
		`(\d+)`:        h.AssertResponseBodyPathIsEqualToFloat64,
		`(true|false)`: h.AssertResponseBodyPathIsEqualToBoolean,
	}
	for pattern, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is %s$`, pattern), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) is %s$`, pattern), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) is %s$`, pattern), f)
	}
	patterns = map[string]any{
		`"([^"]*)"`:    h.AssertNamedHttpRequestResponseBodyPathEqualsTo,
		valueRegex:     h.AssertNamedHttpRequestResponseBodyPathEqualsToValue,
		`(\d+)`:        h.AssertNamedHttpRequestResponseBodyPathIsEqualToFloat64,
		`(true|false)`: h.AssertNamedHttpRequestResponseBodyPathIsEqualToBoolean,
	}
	for pattern, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) for %s is %s$`, httpRequestRegex, pattern), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*)  for %s is %s$`, httpRequestRegex, pattern), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) for %s is %s$`, httpRequestRegex, pattern), f)
	}
	h.ctx.GerkhinContext().Then(`^(?i)response body is file://(.+)$`, h.AssertResponseBodyEqualsToFile)
	h.ctx.GerkhinContext().Then(`^(?i)the response body is file://(.+)$`, h.AssertResponseBodyEqualsToFile)
	h.ctx.GerkhinContext().Then(`^(?i)the Response body is file://(.+)$`, h.AssertResponseBodyEqualsToFile)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body for %s is file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body for %s is file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body for %s is file://(.+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsToFile)

	h.ctx.GerkhinContext().Then(`^(?i)response body is:$`, h.AssertResponseBodyEqualsTo)
	h.ctx.GerkhinContext().Then(`^(?i)the response body is:$`, h.AssertResponseBodyEqualsTo)
	h.ctx.GerkhinContext().Then(`^(?i)the Response body is:$`, h.AssertResponseBodyEqualsTo)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body for %s is:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsTo)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body for %s is:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsTo)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body for %s is:$`, httpRequestRegex), h.AssertNamedHttpRequestResponseBodyEqualsTo)

}

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

func Set(ctx v3.ScenarioContext) {
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
	patterns = map[string]any{
		"contains":    []any{h.AssertResponsePathContains, h.AssertWhileIgnoringCaseThatResponsePathContains},
		"ends with":   []any{h.AssertResponsePathEndsWith, h.AssertWhileIgnoringCaseThatResponsePathEndsWith},
		"starts with": []any{h.AssertResponsePathStartsWith, h.AssertWhileIgnoringCaseThatResponsePathStartsWith},
	}
	valueOpts := []string{`"([^"]*)"$`, valueRegex}
	for k, arr := range patterns {
		for arrayIndex, attr := range arr.([]any) {
			for _, opt := range valueOpts {
				expr := k
				if arrayIndex == 1 {
					expr = "ignoring case " + k
				}
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) %s %s$`, expr, opt), attr)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) %s %s$`, expr, opt), attr)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) %s %s$`, expr, opt), attr)
			}
		}
	}
	patterns = map[string]any{
		"contains":    []any{h.AssertNamedHttpRequestResponsePathContains, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathContains},
		"ends with":   []any{h.AssertNamedHttpRequestResponsePathEndsWith, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathEndsWith},
		"starts with": []any{h.AssertNamedHttpRequestResponsePathStartsWith, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathStartsWith},
	}
	for k, arr := range patterns {
		for arrayIndex, attr := range arr.([]any) {
			for _, opt := range valueOpts {
				expr := k
				if arrayIndex == 1 {
					expr = "ignoring case " + k
				}
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, opt), attr)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, opt), attr)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, opt), attr)
			}
		}
	}
	h.ctx.GerkhinContext().Then(`^(?i)response body \$.(.*) matches pattern "([^"]*)"$`, h.AssertResponsePathMatchesPattern)
	h.ctx.GerkhinContext().Then(`^(?i)the response body \$.(.*) matches pattern "([^"]*)"$`, h.AssertResponsePathMatchesPattern)
	h.ctx.GerkhinContext().Then(`^(?i)the Response body \$.(.*) matches pattern "([^"]*)"$`, h.AssertResponsePathMatchesPattern)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*),for %s, matches pattern "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathMatchesPattern)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, matches pattern "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathMatchesPattern)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, matches pattern "([^"]*)"$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathMatchesPattern)
	patterns = map[string]any{
		"Time":     h.AssertResponsePathIsTime,
		"Date":     h.AssertResponsePathIsDate,
		"DateTime": h.AssertResponsePathIsDateTime,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is %s$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) is %s$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) is %s$`, expr), f)
	}
	patterns = map[string]any{
		"Time":     h.AssertNamedHttpRequestResponsePathIsTime,
		"Date":     h.AssertNamedHttpRequestResponsePathIsDate,
		"DateTime": h.AssertNamedHttpRequestResponsePathIsDateTime,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) for %s is %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) for %s is %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) for %s is %s$`, httpRequestRegex, expr), f)
	}
	patterns = map[string]any{
		"is same":            h.AssertResponsePathIsSame,
		"is after":           h.AssertResponsePathIsAfter,
		"is before":          h.AssertResponsePathIsBefore,
		"is same or after":   h.AssertResponsePathIsSameOrAfter,
		"is before or after": h.AssertResponsePathIsSameOrBefore,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) %s "([^"]*)"$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) %s "([^"]*)"$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) %s "([^"]*)"$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) %s %s$`, expr, valueRegex), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) %s %s$`, expr, valueRegex), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) %s %s$`, expr, valueRegex), f)
	}
	patterns = map[string]any{
		"is same":            h.AssertNamedHttpRequestResponsePathIsSame,
		"is after":           h.AssertNamedHttpRequestResponsePathIsAfter,
		"is before":          h.AssertNamedHttpRequestResponsePathIsBefore,
		"is same or after":   h.AssertNamedHttpRequestResponsePathIsSameOrAfter,
		"is before or after": h.AssertNamedHttpRequestResponsePathIsSameOrBefore,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*), for %s, %s "([^"]*)"$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, %s "([^"]*)"$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, %s "([^"]*)"$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, valueRegex), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, valueRegex), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, %s %s$`, httpRequestRegex, expr, valueRegex), f)
	}
	h.ctx.GerkhinContext().Then(`^(?i)response body \$.(.*) length is (\d+)$`, h.AssertResponsePathLengthIs)
	h.ctx.GerkhinContext().Then(`^(?i)the response body \$.(.*) length is (\d+)$`, h.AssertResponsePathLengthIs)
	h.ctx.GerkhinContext().Then(`^(?i)the Response body \$.(.*) length is (\d+)$`, h.AssertResponsePathLengthIs)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) length for %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathLengthIs)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) length for %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathLengthIs)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) length for %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponsePathLengthIs)
	patterns = map[string]any{
		"lesser":              []any{h.AssertResponsePathIsLesserThan, h.AssertResponsePathIsLesserThanValue},
		"greater than":        []any{h.AssertResponsePathIsGreaterThan, h.AssertResponsePathIsGreaterThanValue},
		"lesser or equal to":  []any{h.AssertResponsePathIsLesserThanOrEqualTo, h.AssertResponsePathIsLesserThanOrEqualToValue},
		"greater or equal to": []any{h.AssertResponsePathIsGreaterThanOrEqualTo, h.AssertResponsePathIsGreaterThanOrEqualToValue},
	}
	for k, opts := range patterns {
		for i, f := range opts.([]any) {
			if i == 0 {
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is %s (\d+)$`, k), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) is %s (\d+)$`, k), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) is %s (\d+)$`, k), f)
			} else {
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is %s %s$`, k, valueRegex), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) is %s %s$`, k, valueRegex), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) is %s %s$`, k, valueRegex), f)
			}
		}
	}
	patterns = map[string]any{
		"lesser":              []any{h.AssertNamedHttpRequestResponsePathIsLesserThan, h.AssertNamedHttpRequestResponsePathIsLesserThanValue},
		"greater than":        []any{h.AssertNamedHttpRequestResponsePathIsGreaterThan, h.AssertNamedHttpRequestResponsePathIsGreaterThanValue},
		"lesser or equal to":  []any{h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualToValue},
		"greater or equal to": []any{h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualToValue},
	}
	for k, opts := range patterns {
		for i, f := range opts.([]any) {
			if i == 0 {
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is %s, for %s, (\d+)$`, httpRequestRegex, k), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, is %s (\d+)$`, httpRequestRegex, k), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, is %s (\d+)$`, httpRequestRegex, k), f)
			} else {
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*), for %s, is %s %s$`, httpRequestRegex, k, valueRegex), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, is %s %s$`, httpRequestRegex, k, valueRegex), f)
				h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, is %s %s$`, httpRequestRegex, k, valueRegex), f)
			}
		}
	}
	patterns = map[string]any{
		`\["[^"]*"(?:,"[^"]*")*\]`: h.AssertResponsePathIsInStringArray,
		`\["\d+"(?:,"\d+")*\]`:     h.AssertResponsePathIsInNumericArray,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*) is part of %s$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*) is part of %s$`, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*) is part of %s$`, expr), f)
	}

	patterns = map[string]any{
		`\["[^"]*"(?:,"[^"]*")*\]`: h.AssertNamedHttpRequestResponsePathIsInStringArray,
		`\["\d+"(?:,"\d+")*\]`:     h.AssertNamedHttpRequestResponsePathIsInNumericArray,
	}
	for expr, f := range patterns {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body \$.(.*), for %s, is part of %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body \$.(.*), for %s, is part of %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body \$.(.*), for %s, is part of %s$`, httpRequestRegex, expr), f)
	}

}

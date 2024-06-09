package v3

import (
	"fmt"
	v3 "github.com/raitonbl/coverup/internal/v3"
)

const serverURLRegex = `(https?://[^\s]+)`
const relativeURIRegex = `/([^/]+(?:/[^/]+)*)`
const valueRegex = `\{\{\s*([a-zA-Z0-9_]+\.)*[a-zA-Z0-9_]+\s*\}\}`
const httpRequestRegex = `\{\{\s*HttpRequest\.[a-zA-Z0-9_]+\s*\}\}`
const entityRegex = `\{\{\s*Entity\.[a-zA-Z0-9_]+\s*\}\}`

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
	onResponseHeaders(h)
	onResponseStatusCode(h)
	onResponseBodySchemaValidation(h)
	onResponseBody(h)
}

func onResponseStatusCode(h *HttpContext) {
	h.ctx.GerkhinContext().Then(`^(?i)response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(`^(?i)the response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(`^(?i)the Response status code is (\d+)$`, h.AssertResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response status code for %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response status code %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response status code %s is (\d+)$`, httpRequestRegex), h.AssertNamedHttpRequestResponseStatusCode)
}

func onResponseHeaders(h *HttpContext) {
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
}

func onResponseBodyPathCompareTo(h *HttpContext) {
	patterns := map[string][]any{
		`"([^"]*)"`:       {h.AssertResponseBodyPathEqualsTo, h.AssertNamedHttpRequestResponseBodyPathEqualsTo},
		valueRegex:        {h.AssertResponseBodyPathEqualsToValue, h.AssertNamedHttpRequestResponseBodyPathEqualsToValue},
		`(-?\d+(\.\d+)?)`: {h.AssertResponseBodyPathIsEqualToFloat64, h.AssertNamedHttpRequestResponseBodyPathIsEqualToFloat64},
		`(true|false)`:    {h.AssertResponseBodyPathIsEqualToBoolean, h.AssertNamedHttpRequestResponseBodyPathIsEqualToBoolean},
	}
	for pattern, opts := range patterns {
		assertResponseBodyPath(h, `is `+pattern, opts[0], opts[1])
	}
	patterns = map[string][]any{
		"contains":    {h.AssertResponsePathContains, h.AssertWhileIgnoringCaseThatResponsePathContains, h.AssertNamedHttpRequestResponsePathContains, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathContains},
		"ends with":   {h.AssertResponsePathEndsWith, h.AssertWhileIgnoringCaseThatResponsePathEndsWith, h.AssertNamedHttpRequestResponsePathEndsWith, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathEndsWith},
		"starts with": {h.AssertResponsePathStartsWith, h.AssertWhileIgnoringCaseThatResponsePathStartsWith, h.AssertNamedHttpRequestResponsePathStartsWith, h.AssertWhileIgnoringCaseThatNamedHttpRequestResponsePathStartsWith},
	}
	valueOpts := []string{`"([^"]*)"$`, valueRegex}
	for k, arr := range patterns {
		for i := 0; i < 2; i++ {
			for _, opt := range valueOpts {
				expr := k
				if i == 1 {
					expr = "ignoring case " + k
				}
				assertResponseBodyPath(h, fmt.Sprintf(`%s %s$`, expr, opt), arr[i], arr[i+2])
			}
		}
	}
}

func onResponseBody(h *HttpContext) {
	assertResponseBody(h, `is:$`, h.AssertResponseBodyEqualsToFile, h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	assertResponseBody(h, `is file://(.+)$`, h.AssertResponseBodyEqualsToFile, h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	onResponseBodyPathCompareTo(h)
	patterns := map[string][]any{}
	assertResponseBodyPath(h, `matches pattern "([^"]*)"$`, h.AssertResponsePathMatchesPattern, h.AssertNamedHttpRequestResponsePathMatchesPattern)
	patterns = map[string][]any{
		"Time":     {h.AssertResponsePathIsTime, h.AssertNamedHttpRequestResponsePathIsTime},
		"Date":     {h.AssertResponsePathIsDate, h.AssertNamedHttpRequestResponsePathIsDate},
		"DateTime": {h.AssertResponsePathIsDateTime, h.AssertNamedHttpRequestResponsePathIsDateTime},
	}
	for expr, arr := range patterns {
		assertResponseBodyPath(h, fmt.Sprintf(`is %s`, expr), arr[0], arr[1])
	}
	patterns = map[string][]any{
		"is same":            {h.AssertResponsePathIsSame, h.AssertNamedHttpRequestResponsePathIsSame},
		"is after":           {h.AssertResponsePathIsAfter, h.AssertNamedHttpRequestResponsePathIsAfter},
		"is before":          {h.AssertResponsePathIsBefore, h.AssertNamedHttpRequestResponsePathIsBefore},
		"is same or after":   {h.AssertResponsePathIsSameOrAfter, h.AssertNamedHttpRequestResponsePathIsSameOrAfter},
		"is before or after": {h.AssertResponsePathIsSameOrBefore, h.AssertNamedHttpRequestResponsePathIsSameOrBefore},
	}
	for expr, arr := range patterns {
		assertResponseBodyPath(h, fmt.Sprintf(`%s "([^"]*)"$`, expr), arr[0], arr[1])
		assertResponseBodyPath(h, fmt.Sprintf(`%s %s$`, expr, valueRegex), arr[0], arr[1])
	}
	assertResponseBodyPath(h, `length is (\d+)`, h.AssertResponsePathLengthIs, h.AssertNamedHttpRequestResponsePathLengthIs)
	patterns = map[string][]any{
		"lesser":              {h.AssertResponsePathIsLesserThan, h.AssertResponsePathIsLesserThanValue, h.AssertNamedHttpRequestResponsePathIsLesserThan, h.AssertNamedHttpRequestResponsePathIsLesserThanValue},
		"greater than":        {h.AssertResponsePathIsGreaterThan, h.AssertResponsePathIsGreaterThanValue, h.AssertNamedHttpRequestResponsePathIsGreaterThan, h.AssertNamedHttpRequestResponsePathIsGreaterThanValue},
		"lesser or equal to":  {h.AssertResponsePathIsLesserThanOrEqualTo, h.AssertResponsePathIsLesserThanOrEqualToValue, h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualToValue},
		"greater or equal to": {h.AssertResponsePathIsGreaterThanOrEqualTo, h.AssertResponsePathIsGreaterThanOrEqualToValue, h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualToValue},
	}
	for k, opts := range patterns {
		assertResponseBodyPath(h, fmt.Sprintf(`is %s (-?\d+(\.\d+)?)$`, k), opts[0], opts[2])
		assertResponseBodyPath(h, fmt.Sprintf(`is %s %s$`, k, valueRegex), opts[1], opts[3])
	}
	patterns = map[string][]any{
		`(\["[^"]*"(?:,"[^"]*")*\])`: {h.AssertResponsePathIsInStringArray, h.AssertNamedHttpRequestResponsePathIsInStringArray},
		`(\["\d+"(?:,"\d+")*\])`:     {h.AssertResponsePathIsInNumericArray, h.AssertNamedHttpRequestResponsePathIsInNumericArray},
	}
	for expr, arr := range patterns {
		assertResponseBodyPath(h, fmt.Sprintf(`is part of %s`, expr), arr[0], arr[1])
	}
}

func onResponseBodySchemaValidation(h *HttpContext) {
	assertResponseBody(h, `respects schema file://(.+)$`, h.AssertResponseIsValidAgainstSchema, h.AssertNamedHttpRequestResponseIsValidAgainstSchema)
}

func assertResponseBodyPath(h *HttpContext, target string, f any, namedF any) {
	assertResponseBody(h, `\$.(.*) `+target, f, namedF)
}

func assertResponseBody(h *HttpContext, target string, f any, namedF any) {
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body %s$`, target), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body %s$`, target), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body %s$`, target), f)
	if namedF != nil {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body, for %s, %s$`, httpRequestRegex, target), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body, for %s, %s$`, httpRequestRegex, target), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body, for %s, %s$`, httpRequestRegex, target), f)
	}
}

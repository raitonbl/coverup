package v3

import (
	"fmt"
	"strings"
)

const serverURLRegex = `(https?://[^\s]+)`
const relativeURIRegex = `/([^/]+(?:/[^/]+)*)`
const valueRegex = `\{\{\s*([a-zA-Z0-9_]+\.)*[a-zA-Z0-9_]+\s*\}\}`
const httpRequestRegex = `\{\{\s*HttpRequest\.[a-zA-Z0-9_]+\s*\}\}`
const entityRegex = `\{\{\s*Entity\.[a-zA-Z0-9_]+\s*\}\}`

func On(ctx ScenarioContext) {
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
	setRequestLinePart(h, `method is ([^"]*)$`, h.WithMethod)
	// Define headers
	setRequestLinePart(h, `headers:$`, h.WithHeaders)
	setRequestLinePart(h, fmt.Sprintf(`header (.*) is "%s"$`, valueRegex), h.WithHeader)
	setRequestLinePart(h, `header (.*) is "([^"]*)"$`, h.WithHeader)
	setRequestLinePart(h, `accept is "([^"]*)"$`, h.WithAcceptHeader)
	setRequestLinePart(h, `content-type is "([^"]*)"$`, h.WithContentTypeHeader)
	// Server URL & Path
	setRequestLinePart(h, `path is http://(.+)$`, h.WithHttpPath)
	setRequestLinePart(h, `path is https://(.+)$`, h.WithHttpsPath)
	setRequestLinePart(h, fmt.Sprintf(`path is %s$`, valueRegex), h.WithPath)
	setRequestLinePart(h, fmt.Sprintf(`path is %s$`, relativeURIRegex), h.WithPath)
	setRequestLinePart(h, fmt.Sprintf(`server url is %s$`, serverURLRegex), h.WithServerURL)
	//Body
	setRequestLinePart(h, `body is:$`, h.WithBody)
	setRequestLinePart(h, `body is file://(.+)$`, h.WithBodyFileURI)
	//Form
	setRequestLinePart(h, `form enctype is ([^"]*)$`, h.WithFormEncType)
	setRequestLinePart(h, `form attribute "([a-zA-Z_]+)" is "([^"]+)"$`, h.WithFormAttribute)
	setRequestLinePart(h, fmt.Sprintf(`form attribute "([a-zA-Z_]+)" is "%s"$`, valueRegex), h.WithFormAttribute)
	setRequestLinePart(h, `form attribute "([a-zA-Z_]+)" is file://(.+)$`, h.WithFormFile)
	setRequestLinePart(h, fmt.Sprintf(`form attribute "([a-zA-Z_]+)" is file://%s`, valueRegex), h.WithFormFile)
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

func setRequestLinePart(h *HttpContext, expr string, f any) {
	h.ctx.GerkhinContext().Step(`^(?i)`+expr, f)
	h.ctx.GerkhinContext().Step(`^(?i)the `+expr, f)
	h.ctx.GerkhinContext().Step(`^(?i)the `+strings.ToUpper(string(expr[0]))+expr[1:], f)
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
	vars := map[string]HandlerFactory{
		`"([^"]*)"`:       HandlerFactory(createResponseBodyPathEqualTo),
		valueRegex:        HandlerFactory(createResponseBodyPathEqualTo),
		`(-?\d+(\.\d+)?)`: HandlerFactory(createResponseBodyPathEqualToFloat64),
		`(true|false)`:    HandlerFactory(createResponseBodyPathEqualToBoolean),
	}
	verbs := []string{"is", "isn't"}
	for pattern, f := range vars {
		for _, verb := range verbs {
			isTrue := verb == verbs[0]
			assertResponseBodyPath(h, verb+` `+pattern, f(h, Opts{isAffirmation: isTrue, isAliasAware: false}), f(h, Opts{isAffirmation: isTrue, isAliasAware: true}))
		}
	}
	vars = map[string]HandlerFactory{
		"contains":    createResponseBodyPathContains,
		"ends with":   createResponseBodyPathEndsWith,
		"starts with": createResponseBodyPathStartsWith,
	}
	n := map[string]string{
		"contains":    "contain",
		"ends with":   "end with",
		"starts with": "start with",
	}
	valueOpts := []string{`"([^"]*)"$`, valueRegex}
	verbs = []string{"", "doesn't"}
	for k, f := range vars {
		for i := 0; i < 2; i++ {
			for _, opt := range valueOpts {
				for _, verb := range verbs {
					pattern := ""
					if i == 1 {
						pattern = "ignoring case"
					}
					if verb == verbs[1] {
						pattern += " " + verb + " " + n[k]
					} else if pattern == "" {
						pattern = k
					} else {
						pattern += " " + k
					}
					isTrue := verb == verbs[0]
					assertResponseBodyPath(h, pattern+` `+opt, f(h, Opts{isAffirmation: isTrue, isAliasAware: false, ignoreCase: i == 1, interpolateValue: true}),
						f(h, Opts{isAffirmation: isTrue, isAliasAware: true, ignoreCase: i == 1, interpolateValue: true}))
				}
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

func assertResponseBodyPath(h *HttpContext, expr string, f any, namedF any) {
	assertResponseBody(h, fmt.Sprintf(`\$(\S+) %s`, expr), f, namedF)
}

func assertResponseBody(h *HttpContext, expr string, f any, namedF any) {
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body %s$`, expr), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body %s$`, expr), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body %s$`, expr), f)
	if namedF != nil {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body, for %s, %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body, for %s, %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body, for %s, %s$`, httpRequestRegex, expr), f)
	}
}

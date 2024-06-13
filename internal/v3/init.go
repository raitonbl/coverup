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

func onResponseBody(h *HttpContext) {
	//	setRequestBodyStepDefinition(h, `is:$`, h.AssertResponseBodyEqualsToFile, h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	//	setRequestBodyStepDefinition(h, `is file://(.+)$`, h.AssertResponseBodyEqualsToFile, h.AssertNamedHttpRequestResponseBodyEqualsToFile)
	params := map[string]HandlerFactory{
		":":            newResponseBodyIsEqualToHandler,
		" file://(.*)": newResponseBodyIsEqualToFileHandler,
	}
	verbs := []string{"is", "isn't"}
	for expr, f := range params {
		for _, verb := range verbs {
			isAffirmation := verb == verbs[0]
			setRequestBodyStepDefinition(h, verb+expr,
				f(h, HandlerOpts{isAffirmationExpected: isAffirmation, isAliasedFunction: false}),
				f(h, HandlerOpts{isAffirmationExpected: isAffirmation, isAliasedFunction: true}))
		}
	}
	onJsonPathCompareTo(h)
	patterns := map[string][]any{}
	patterns = map[string][]any{
		"Time":     {h.AssertResponsePathIsTime, h.AssertNamedHttpRequestResponsePathIsTime},
		"Date":     {h.AssertResponsePathIsDate, h.AssertNamedHttpRequestResponsePathIsDate},
		"DateTime": {h.AssertResponsePathIsDateTime, h.AssertNamedHttpRequestResponsePathIsDateTime},
	}
	for expr, arr := range patterns {
		setJsonPathStepDefinition(h, fmt.Sprintf(`is %s`, expr), arr[0], arr[1])
	}
	patterns = map[string][]any{
		"is same":            {h.AssertResponsePathIsSame, h.AssertNamedHttpRequestResponsePathIsSame},
		"is after":           {h.AssertResponsePathIsAfter, h.AssertNamedHttpRequestResponsePathIsAfter},
		"is before":          {h.AssertResponsePathIsBefore, h.AssertNamedHttpRequestResponsePathIsBefore},
		"is same or after":   {h.AssertResponsePathIsSameOrAfter, h.AssertNamedHttpRequestResponsePathIsSameOrAfter},
		"is before or after": {h.AssertResponsePathIsSameOrBefore, h.AssertNamedHttpRequestResponsePathIsSameOrBefore},
	}
	for expr, arr := range patterns {
		setJsonPathStepDefinition(h, fmt.Sprintf(`%s "([^"]*)"$`, expr), arr[0], arr[1])
		setJsonPathStepDefinition(h, fmt.Sprintf(`%s %s$`, expr, valueRegex), arr[0], arr[1])
	}
	setJsonPathStepDefinition(h, `length is (\d+)`, h.AssertResponsePathLengthIs, h.AssertNamedHttpRequestResponsePathLengthIs)
	patterns = map[string][]any{
		"lesser":              {h.AssertResponsePathIsLesserThan, h.AssertResponsePathIsLesserThanValue, h.AssertNamedHttpRequestResponsePathIsLesserThan, h.AssertNamedHttpRequestResponsePathIsLesserThanValue},
		"greater than":        {h.AssertResponsePathIsGreaterThan, h.AssertResponsePathIsGreaterThanValue, h.AssertNamedHttpRequestResponsePathIsGreaterThan, h.AssertNamedHttpRequestResponsePathIsGreaterThanValue},
		"lesser or equal to":  {h.AssertResponsePathIsLesserThanOrEqualTo, h.AssertResponsePathIsLesserThanOrEqualToValue, h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsLesserThanOrEqualToValue},
		"greater or equal to": {h.AssertResponsePathIsGreaterThanOrEqualTo, h.AssertResponsePathIsGreaterThanOrEqualToValue, h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualTo, h.AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualToValue},
	}
	for k, opts := range patterns {
		setJsonPathStepDefinition(h, fmt.Sprintf(`is %s (-?\d+(\.\d+)?)$`, k), opts[0], opts[2])
		setJsonPathStepDefinition(h, fmt.Sprintf(`is %s %s$`, k, valueRegex), opts[1], opts[3])
	}
	patterns = map[string][]any{
		`(\["[^"]*"(?:,"[^"]*")*\])`: {h.AssertResponsePathIsInStringArray, h.AssertNamedHttpRequestResponsePathIsInStringArray},
		`(\["\d+"(?:,"\d+")*\])`:     {h.AssertResponsePathIsInNumericArray, h.AssertNamedHttpRequestResponsePathIsInNumericArray},
	}
	for expr, arr := range patterns {
		setJsonPathStepDefinition(h, fmt.Sprintf(`is part of %s`, expr), arr[0], arr[1])
	}
}

func onResponseBodySchemaValidation(h *HttpContext) {
	schemes := []URIScheme{
		fileUriScheme,
		httpUriScheme,
		httpsUriScheme,
	}
	verbs := []string{"", "doesn't"}
	for _, scheme := range schemes {
		for _, verb := range verbs {
			prefix := verb
			isAffirmation := verb == verbs[0]
			if verb == verbs[1] {
				prefix += " "
			}
			setRequestBodyStepDefinition(h, fmt.Sprintf(`%srespects schema %s://(.+)$`, prefix, scheme),
				newJsonSchemaValidator(h, HandlerOpts{isAffirmationExpected: isAffirmation, isAliasedFunction: false, scheme: scheme}),
				HandlerOpts{isAffirmationExpected: isAffirmation, isAliasedFunction: true, scheme: scheme})
		}
	}
}

func onJsonPathCompareTo(h *HttpContext) {
	doOnJsonPathCompareTo(h, []string{"is", "isn't"}, map[string]HandlerFactory{
		`"([^"]*)"`:       newJsonPathEqualsTo,
		valueRegex:        newJsonPathEqualsTo,
		`(-?\d+(\.\d+)?)`: newJsonPathEqualsToFloat64,
		`(true|false)`:    newJsonPathEqualsToBooleanHandler,
	}, []HandlerOpts{
		{isAffirmationExpected: true, isAliasedFunction: false},
		{isAffirmationExpected: true, isAliasedFunction: true},
		{isAffirmationExpected: false, isAliasedFunction: false},
		{isAffirmationExpected: false, isAliasedFunction: true},
	})
	doOnJsonPathStringOperation(h, []string{"", "doesn't"})
}

func doOnJsonPathStringOperation(h *HttpContext, verbs []string) {
	patterns := map[string]HandlerFactory{
		"contains":        newJsonPathContainsHandler,
		"ends with":       newJsonPathEndsWithHandler,
		"starts with":     newJsonPathStartsWithHandler,
		"matches pattern": newJsonPathPatternHandler,
	}
	negations := map[string]string{
		"contains":        "contain",
		"ends with":       "end with",
		"starts with":     "start with",
		"matches pattern": "match pattern",
	}
	valueOpts := []string{`"([^"]*)"$`, valueRegex}
	for pattern, factory := range patterns {
		for _, valueOpt := range valueOpts {
			for i := 0; i < 2; i++ {
				ignoreCase := i == 1
				for _, verb := range verbs {
					opt := HandlerOpts{isAffirmationExpected: verb == verbs[0], isAliasedFunction: false, ignoreCaseIfApplicable: ignoreCase, attemptValueResolution: true}
					assertionPattern := pattern
					if ignoreCase {
						assertionPattern = "ignoring case " + assertionPattern
					}
					if verb == "doesn't" {
						assertionPattern = verb + " " + negations[pattern]
					}
					setJsonPathStepDefinition(h, assertionPattern+" "+valueOpt, factory(h, HandlerOpts(opt)), factory(h, HandlerOpts(opt)))
				}
			}
		}
	}
}

func doOnJsonPathCompareTo(h *HttpContext, verbs []string, patterns map[string]HandlerFactory, opts []HandlerOpts) {
	for pattern, factory := range patterns {
		for _, verb := range verbs {
			for _, opt := range opts {
				isAffirmation := verb == verbs[0]
				assertionPattern := verb + " " + pattern
				if opt.ignoreCaseIfApplicable {
					assertionPattern = "ignoring case " + assertionPattern
				}
				target := HandlerOpts{
					isAffirmationExpected:  isAffirmation,
					isAliasedFunction:      opt.isAliasedFunction,
					ignoreCaseIfApplicable: opt.ignoreCaseIfApplicable,
					attemptValueResolution: opt.attemptValueResolution,
				}
				setJsonPathStepDefinition(h, assertionPattern, factory(h, target), factory(h, target))
			}
		}
	}
}

func setJsonPathStepDefinition(h *HttpContext, expr string, f any, namedF any) {
	setRequestBodyStepDefinition(h, fmt.Sprintf(`\$(\S+) %s`, expr), f, namedF)
}

func setRequestBodyStepDefinition(h *HttpContext, expr string, f any, namedF any) {
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body %s$`, expr), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body %s$`, expr), f)
	h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body %s$`, expr), f)
	if namedF != nil {
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)response body, for %s, %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the response body, for %s, %s$`, httpRequestRegex, expr), f)
		h.ctx.GerkhinContext().Then(fmt.Sprintf(`^(?i)the Response body, for %s, %s$`, httpRequestRegex, expr), f)
	}
}

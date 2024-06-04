package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

const namedHttpRequestRegex = `\{\{HttpRequest\.(\w+)\}\}`

var schemas = []string{"file", "http", "https"}
var methods = []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}

func Configure(ctx context.Context, scenarioContext *godog.ScenarioContext) error {
	b := context.New(ctx)
	// Request
	scenarioContext.Given(`^a HttpRequest$`, CreateHttpRequest(b))
	scenarioContext.Given(`^a HttpRequest named (.+)$`, CreateHttpRequestWithAlias(b))
	scenarioContext.Step(`^the headers:$`, SetRequestHeaders(b))
	scenarioContext.Step(`^(?i)body is:$`, SetRequestBody(b))
	scenarioContext.When(`^(?i)submitting HttpRequest$`, SubmitsHttpRequest(b))
	scenarioContext.When(fmt.Sprintf(`^(?i)submitting %s`, namedHttpRequestRegex),
		SubmitsHttpRequestWhenAlias(b))
	for _, method := range methods {
		scenarioContext.Step(`^(?i)operation `+method+` (\/[^\s]*|https?://[^\s]+)$`, SetRequestOperation(b, method))
	}
	// Common
	withURI(b, scenarioContext)
	// Response
	scenarioContext.Then(`^the response status code is (\d+)$`, AssertHttpResponseStatusCode(b))
	scenarioContext.Then(fmt.Sprintf(`^the %s response status code is (\d+)$`, namedHttpRequestRegex),
		AssertHttpResponseStatusCodeWhenAlias(b))

	return nil
}

func withURI(b *context.Builder, scenarioContext *godog.ScenarioContext) {
	for _, schemaType := range schemas {
		f := SetRequestBodyFromURI(b, schemaType)
		// Request
		scenarioContext.Step(fmt.Sprintf(`^body is %s://(.+)$`, schemaType), f)
		scenarioContext.Step(fmt.Sprintf(`^the body is %s://(.+)$`, schemaType), f)
		// Response
		scenarioContext.Then(fmt.Sprintf(`^the response body complies with schema %s://(.+)$`, schemaType),
			AssertHttpResponseBodySchemaOnURI(b, schemaType))
		scenarioContext.Then(fmt.Sprintf(`^the %s response body complies with schema %s://(.+)$`, namedHttpRequestRegex, schemaType),
			AssertHttpResponseBodySchemaOnURIWhenAlias(b, schemaType))
	}
}

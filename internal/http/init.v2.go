package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

const namedHttpRequestRegex = "{{\\sHttpRequest.(.+)\\s}}"

var schemas = []string{"file", "http", "https"}
var methods = []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}

func Configure(ctx context.Context, scenarioContext *godog.ScenarioContext) error {
	b := context.New(ctx)
	// Request
	scenarioContext.Given(`^a HttpRequest$`, CreateHttpRequest(b))
	scenarioContext.Given(`^a HttpRequest named (.+)$`, CreateHttpRequestWithAlias(b))
	scenarioContext.Given(`^the headers:$`, SetRequestHeaders(b))
	scenarioContext.Given(`^body is:$`, SetRequestBody(b))
	scenarioContext.Then(`^submits HttpRequest$`, SubmitsHttpRequest(b))
	scenarioContext.Then(fmt.Sprintf(`^submits %s`, namedHttpRequestRegex),
		SubmitsHttpRequestWhenAlias(b))
	for _, method := range methods {
		scenarioContext.When(fmt.Sprintf(`^operation %s "([^"]*)"$`, method), SetRequestOperation(b, method))
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
		scenarioContext.Given(fmt.Sprintf(`^body is %s://(.+)$`, schemaType), f)
		scenarioContext.Given(fmt.Sprintf(`^the body is %s://(.+)$`, schemaType), f)
		// Response
		scenarioContext.Then(fmt.Sprintf(`^the response body complies with schema %s://(.+)$`, schemaType),
			AssertHttpResponseBodySchemaOnURI(b, schemaType))
		scenarioContext.Then(fmt.Sprintf(`^the %s response body complies with schema %s://(.+)$`, namedHttpRequestRegex, schemaType),
			AssertHttpResponseBodySchemaOnURIWhenAlias(b, schemaType))
	}
}

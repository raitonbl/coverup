package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	http "github.com/raitonbl/coverup/internal/http"
)

func NewV2(ctx context.Context) func(*godog.ScenarioContext) {
	return func(scenarioContext *godog.ScenarioContext) {
		b := context.New(ctx)
		// Given Params
		scenarioContext.Given(`^a HttpRequest$`, http.CreateHttpRequest(b))
		scenarioContext.Given(`^a HttpRequest <(.*)>$`, http.CreateHttpRequestWithAlias(b))
		scenarioContext.Given(`^the headers:$`, http.SetRequestHeaders(b))
		// When
		methods := []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}
		for _, method := range methods {
			f := http.SetRequestOperation(b, method)
			scenarioContext.When(fmt.Sprintf(`^Operation %s "([^"]*)"$`, method), f)
		}
		scenarioContext.Given(`^body is:$`, http.SetRequestBody(b))
		array := []string{"file", "http", "https"}
		for _, schemaType := range array {
			f := http.SetRequestBodyFromURI(b, schemaType)
			scenarioContext.Given(fmt.Sprintf(`^body is %s://(.+)$`, schemaType), f)
			scenarioContext.Given(fmt.Sprintf(`^the body is %s://(.+)$`, schemaType), f)
		}
		thenHttp(scenarioContext, `^the response statusCode is (\d+)`,
			http.AssertHttpResponseStatusCode(b), http.AssertHttpResponseStatusCodeWhenAlias(b))
	}
}

func thenHttp(s *godog.ScenarioContext, expr string, fn, aliasFn any) {
	if fn != nil {
		s.Then(fmt.Errorf(`%s$`, expr), fn)
	}
	if aliasFn != nil {
		s.Then(fmt.Errorf(`%s <(.*)>$`, expr), aliasFn)
	}
}

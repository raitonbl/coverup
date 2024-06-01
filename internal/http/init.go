package internal

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

func WithHttpRequest(context context.Context) func(*godog.ScenarioContext) {
	return func(s *godog.ScenarioContext) {
		operation := &Operation{Context: context}
		// Given Params
		s.Step(`^the body:$`, operation.withRequestBody)
		s.Given(`^a HttpRequest$`, operation.withHttpRequest)
		s.Step(`^the headers:$`, operation.withRequestHeaders)
		// When
		s.When(`^GET "([^"]*)"$`, operation.withGetMethod)
		s.When(`^PUT "([^"]*)"$`, operation.withPutMethod)
		s.When(`^POST "([^"]*)"$`, operation.withPostMethod)
		s.When(`^PATCH "([^"]*)"$`, operation.withPatchMethod)
		s.When(`^DELETE "([^"]*)"$`, operation.withDeleteMethod)
		// Then
		s.Then(`^the response body:$`, operation.withHttpResponseBody)
		s.Then(`^the response headers:$`, operation.withHttpResponseHeaders)
		s.Then(`^the response statusCode is (\d+)$`, operation.withStatusCode)
		s.Then(`^the response header ([^ ]+) is (.+)$`, operation.withHttpResponseHeader)
		s.Then(`^the response body uri is: ([^"]*)$`, operation.withResponseBodyURI)
	}

}

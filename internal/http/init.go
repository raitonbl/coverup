package internal

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

func WithHttpRequest(context context.Context) func(*godog.ScenarioContext) {
	return func(s *godog.ScenarioContext) {
		operation := &Operation{Context: context}
		// Given Params
		s.Given(`^the body:$`, operation.withRequestBody)
		s.Given(`^a HttpRequest$`, operation.withHttpRequest)
		s.Given(`^the headers:$`, operation.withRequestHeaders)
		// When
		s.When(`^GET "([^"]*)"$`, operation.withGetMethod)
		s.When(`^PUT "([^"]*)"$`, operation.withPutMethod)
		s.When(`^POST "([^"]*)"$`, operation.withPostMethod)
		s.When(`^PATCH "([^"]*)"$`, operation.withPatchMethod)
		s.When(`^DELETE "([^"]*)"$`, operation.withDeleteMethod)
		// Then:
		s.Then(`^the response body:$`, operation.withHttpResponseBody)
		s.Then(`^the response headers:$`, operation.withHttpResponseHeaders)
		s.Then(`^the response statusCode is (\d+)$`, operation.withStatusCode)
		s.Then(`^the response header ([^ ]+) is (.+)$`, operation.withHttpResponseHeader)
		s.Then(`^the response body uri is: ([^"]*)$`, operation.withResponseBodyURI)
		// Then: Criteria > Body
		s.Then(`^the \$(body\..*) matches "([^"]*)"$`, operation.withBodyPathMatches)
		s.Then(`^the \$(body\..*) ends with "([^"]*)"$`, operation.withBodyPathEndsWith)
		s.Then(`^the \$(body\..*) is equal to "([^"]*)"$`, operation.withBodyPathEqualTo)
		s.Then(`^the \$(body\..*) is lesser than (\d+)$`, operation.withBodyPathLesserThan)
		s.Then(`^the \$(body\..*) starts with "([^"]*)"$`, operation.withBodyPathStartsWith)
		s.Then(`^the \$(body\..*) is greater than (\d+)$`, operation.withBodyPathGreaterThan)
		s.Then(`^the \$(body\..*) is lesser or equal to (\d+)$`, operation.withBodyPathLesserOrEqualTo)
		s.Then(`^the \$(body\..*) is greater or equal to (\d+)$`, operation.withBodyPathGreaterOrEqualTo)
		// Then: Criteria > Headers
		s.Then(`^the \$(headers\..*) ends with "([^"]*)"$`, operation.withHeaderEndsWith)
		s.Then(`^the \$(headers\..*) is equal to "([^"]*)"$`, operation.withHeaderEqualTo)
		s.Then(`^the \$(headers\..*) is lesser than (\d+)$`, operation.withHeaderLesserThan)
		s.Then(`^the \$(headers\..*) starts with "([^"]*)"$`, operation.withHeaderStartsWith)
		s.Then(`^the \$(headers\..*) is greater than (\d+)$`, operation.withHeaderGreaterThan)
		s.Then(`^the \$(headers\..*) is lesser or equal to (\d+)$`, operation.withHeaderLesserOrEqualTo)
		s.Then(`^the \$(headers\..*) is greater or equal to (\d+)$`, operation.withHeaderGreaterOrEqualTo)

	}

}

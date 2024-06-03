package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	http "github.com/raitonbl/coverup/internal/http"
)

func NewV2(ctx context.Context) func(*godog.ScenarioContext) {
	return func(s *godog.ScenarioContext) {
		b := context.New(ctx)
		// Given Params
		s.Given(`^a HttpRequest$`, http.CreateHttpRequest(b))
		s.Given(`^a HttpRequest <(.*)>$`, http.CreateHttpRequestWithAlias(b))
		s.Given(`^the headers:$`, http.SetRequestHeaders(b))
		// When
		methods := []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}
		for _, method := range methods {
			f := http.SetRequestOperation(b, method)
			s.When(fmt.Sprintf(`^Operation %s "([^"]*)"$`, method), f)
		}
		s.Given(`^body is:$`, http.SetRequestBody(b))
		array := []string{"file", "http", "https"}
		for _, schemaType := range array {
			f := http.SetRequestBodyFromURI(b, schemaType)
			s.Given(fmt.Sprintf(`^body is %s://(.+)$`, schemaType), f)
			s.Given(fmt.Sprintf(`^the body is %s://(.+)$`, schemaType), f)

		}
	}
}

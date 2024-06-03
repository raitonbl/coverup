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
		s.Given(`^body is file://(.+)$`, http.SetRequestBodyFromURI(b, "file"))
		s.Given(`^body is http://(.+)$`, http.SetRequestBodyFromURI(b, "http"))
		s.Given(`^body is https://(.+)$`, http.SetRequestBodyFromURI(b, "https"))
		s.Given(`^the body is:$`, http.SetRequestBody(b))
		s.Given(`^the body is file://(.+)$`, http.SetRequestBodyFromURI(b, "file"))
		s.Given(`^the body is http://(.+)$`, http.SetRequestBodyFromURI(b, "http"))
		s.Given(`^the body is https://(.+)$`, http.SetRequestBodyFromURI(b, "https"))
	}
}

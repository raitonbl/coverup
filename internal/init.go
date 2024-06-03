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
		s.Given(`^a HttpRequest <(.*)>$`, http.CreateHttpRequestAndAlias(b))
		s.Given(`^the headers:$`, http.CreateHttpRequestHeaders(b))
		// When
		methods := []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}
		for _, method := range methods {
			f := http.CreateHttpRequestOperation(b, method)
			s.When(fmt.Sprintf(`^Operation %s "([^"]*)"$`, method), f)
		}
	}
}

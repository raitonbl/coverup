package http

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	"strings"
)

type ThenHttpResponseStepFactory struct {
}

func (instance *ThenHttpResponseStepFactory) New(ctx api.StepDefinitionContext) {
	instance.thenHeader(ctx)
	instance.thenHeaders(ctx)
	instance.thenStatusCode(ctx)
	instance.thenResponseBody(ctx)
}

func (instance *ThenHttpResponseStepFactory) thenStatusCode(ctx api.StepDefinitionContext) {

}

func (instance *ThenHttpResponseStepFactory) thenHeader(ctx api.StepDefinitionContext) {

}

func (instance *ThenHttpResponseStepFactory) thenHeaders(ctx api.StepDefinitionContext) {
	instance.thenHeadersContains(ctx)
	instance.thenHeadersEqualsTo(ctx)
}

func (instance *ThenHttpResponseStepFactory) thenHeadersContains(ctx api.StepDefinitionContext) {
	instance.doThenHeaders(ctx, func(fromResponse, definition map[string]string) bool {
		for k, v := range fromResponse {
			if definition[k] != v {
				return false
			}
		}
		return true
	}, [][]string{
		{"contains", "has exact headers"},
		{"doesn't contain", "doesn't have exact headers"},
	})
}

func (instance *ThenHttpResponseStepFactory) thenHeadersEqualsTo(ctx api.StepDefinitionContext) {
	instance.doThenHeaders(ctx, func(fromResponse, definition map[string]string) bool {
		if len(fromResponse) != len(definition) {
			return false
		}
		for k, v := range fromResponse {
			if definition[k] != v {
				return false
			}
		}
		return true
	}, [][]string{
		{"is", "has exact headers"},
		{"isn't", "doesn't have exact headers"},
	})
}

func (instance *ThenHttpResponseStepFactory) doThenHeaders(ctx api.StepDefinitionContext, predicate func(fromResponse, definition map[string]string) bool, verbs [][]string) {
	for _, entry := range verbs {
		verb := entry[0]
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a specific %s response %s", ComponentType, entry[1]),
			Options:     make([]api.Option, 0),
		}
		opts := []string{"", httpRequestRegex}
		for _, opt := range opts {
			var phrases []string
			format := fmt.Sprintf(`headers %s:$`, verb)
			if opt == opts[0] {
				phrases = createResponseLinePart(format)
			} else {
				phrases = createAliasedResponseLinePart(format)
			}
			f := func(c api.ScenarioContext) any {
				cfg := FactoryOpts[any]{
					ResolveValueBeforeAssertion: true,
					AssertAlias:                 opt == opts[1],
					AssertTrue:                  verb == verbs[0][0],
				}
				return instance.createThenHeaders(c, cfg, predicate)
			}
			for _, phrase := range phrases {
				step.Options = append(step.Options, api.Option{
					Regexp:         phrase,
					Description:    step.Description,
					HandlerFactory: f,
				})
			}
		}
		ctx.Then(step)
	}
}

func (instance *ThenHttpResponseStepFactory) createThenHeaders(c api.ScenarioContext, cfg FactoryOpts[any], predicate func(fromResponse, definition map[string]string) bool) any {
	f := func(alias string, table *godog.Table) error {
		res, err := instance.getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		definition := make(map[string]string)
		for _, row := range table.Rows {
			k := row.Cells[0].Value
			v := row.Cells[1].Value
			if cfg.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(v)
				if prob != nil {
					return prob
				}
				definition[k] = fmt.Sprintf("%v", valueOf)
			}
		}
		r := predicate(res.headers, definition)
		if cfg.AssertTrue == r {
			return nil
		}
		sb := strings.Builder{}
		for k, v := range res.headers {
			sb.WriteString(k + "=" + v + "\n")
		}
		return fmt.Errorf("response headers:\n%s", sb.String())
	}
	if cfg.AssertAlias {
		return f
	}
	return func(table *godog.Table) error {
		return f("", table)
	}
}

func (instance *ThenHttpResponseStepFactory) thenResponseBody(ctx api.StepDefinitionContext) {

}

func (instance *ThenHttpResponseStepFactory) getHttpResponse(c api.ScenarioContext, alias string) (*HttpResponse, error) {
	component, err := c.GetGivenComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	req, isHttpRequest := component.(*HttpRequest)
	if !isHttpRequest {
		return nil, fmt.Errorf("cannot retrieve %s from context", ComponentType)
	}
	if req.response == nil {
		return nil, fmt.Errorf("%s must be submitted before using the response", ComponentType)
	}
	return req.response, nil
}

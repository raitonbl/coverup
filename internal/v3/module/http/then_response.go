package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	pkgHttp "github.com/raitonbl/coverup/pkg/http"
	"github.com/thoas/go-funk"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type ThenHttpResponseStepFactory struct {
}

func (instance *ThenHttpResponseStepFactory) New(ctx api.StepDefinitionContext) {
	instance.thenHeaders(ctx)
	instance.thenStatusCode(ctx)
	instance.thenResponseBody(ctx)
}

func (instance *ThenHttpResponseStepFactory) thenStatusCode(ctx api.StepDefinitionContext) {
	verbs := []string{
		"is",
		"isn't",
	}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response status status %s equal to a specific http status code", ComponentType, verb),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		for _, opt := range opts {
			var phrases []string
			format := fmt.Sprintf(`status code %s (\d+)$`, verb)
			if opt == opts[0] {
				phrases = createResponseLinePart(format)
			} else {
				phrases = createAliasedResponseLinePart(format)
			}
			f := func(c api.ScenarioContext) any {
				return instance.createThenStatusCode(c, FactoryOpts[any]{
					ResolveValueBeforeAssertion: true,
					AssertAlias:                 opt == opts[1],
					AssertTrue:                  verb == verbs[0],
				})
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

func (instance *ThenHttpResponseStepFactory) createThenStatusCode(c api.ScenarioContext, opts FactoryOpts[any]) any {
	f := func(alias string, statusCode float64) error {
		req, err := instance.getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		if opts.AssertTrue {
			if req.statusCode == statusCode {
				return nil
			}
			return fmt.Errorf("status code should be %v but got %v", statusCode, req.statusCode)
		}
		if req.statusCode != statusCode {
			return nil
		}
		return fmt.Errorf("status code mustn't be %v, yet got %v", statusCode, req.statusCode)
	}
	if !opts.AssertAlias {
		return func(statusCode float64) error {
			return f("", statusCode)
		}
	}
	return f
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
	instance.thenBodyEqualTo(ctx)
	instance.thenBodyEqualFile(ctx)
	instance.thenBodyRespectsJsonSchema(ctx)
}

func (instance *ThenHttpResponseStepFactory) thenBodyEqualTo(ctx api.StepDefinitionContext) {
	instance.thenBodyEqualToSrc(ctx, "body", ":", false)
}

func (instance *ThenHttpResponseStepFactory) thenBodyEqualFile(ctx api.StepDefinitionContext) {
	instance.thenBodyEqualToSrc(ctx, "file", " file://(.*)", true)
}

func (instance *ThenHttpResponseStepFactory) thenBodyEqualToSrc(ctx api.StepDefinitionContext, targetType, regex string, isFileURI bool) {
	verbs := []string{
		"is",
		"isn't",
	}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response body %s equal to a specific %s", ComponentType, verb, targetType),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		for _, opt := range opts {
			var phrases []string
			format := fmt.Sprintf(`body %s%s`, verb, regex)
			if opt == opts[0] {
				phrases = createResponseLinePart(format)
			} else {
				phrases = createAliasedResponseLinePart(format)
			}
			f := func(c api.ScenarioContext) any {
				return instance.createThenBodyEqualsToSrc(c, isFileURI, FactoryOpts[any]{
					ResolveValueBeforeAssertion: true,
					AssertAlias:                 opt == opts[1],
					AssertTrue:                  verb == verbs[0],
				})
			}
			for _, phrase := range phrases {
				step.Options = append(step.Options, api.Option{
					HandlerFactory: f,
					Regexp:         phrase,
					Description:    step.Description,
				})
			}
		}
		ctx.Then(step)
	}
}

func (instance *ThenHttpResponseStepFactory) createThenBodyEqualsToSrc(c api.ScenarioContext, isFileURI bool, opts FactoryOpts[any]) any {
	h := func(alias string, binary []byte) error {
		res, err := instance.getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		var predicate func() (bool, error)
		if res.headers["content-type"] == "application/json" {
			predicate = func() (bool, error) {
				fromResponse := map[string]any{}
				if prob := json.Unmarshal(res.body, &fromResponse); prob != nil {
					return false, prob
				}
				fromValue := map[string]any{}
				if prob := json.Unmarshal(binary, &fromValue); prob != nil {
					return false, prob
				}
				return reflect.DeepEqual(fromResponse, fromValue), nil
			}
		} else {
			predicate = func() (bool, error) {
				if string(res.body) != string(binary) {
					return false, nil
				}
				return true, nil
			}
		}
		r, err := predicate()
		if err != nil {
			return err
		}
		if opts.AssertTrue == r {
			return nil
		}

		if opts.AssertTrue {
			return fmt.Errorf("response isn't equal to expectation.\n%s", string(res.body))
		}
		return fmt.Errorf("response  shouldn't match expectation")
	}
	if isFileURI {
		f := func(alias, value string) error {
			binary, err := c.GetFS().ReadFile(value)
			if err != nil {
				return err
			}
			return h(alias, binary)
		}
		if opts.AssertAlias {
			return func(value string) error {
				return f("", value)
			}
		}
		return f
	}

	f := func(alias string, value *godog.DocString) error {
		return h(alias, []byte(value.Content))
	}
	if opts.AssertAlias {
		return func(value *godog.DocString) error {
			return f("", value)
		}
	}
	return f
}

func (instance *ThenHttpResponseStepFactory) thenBodyRespectsJsonSchema(ctx api.StepDefinitionContext) {
	verbs := []string{
		"respects",
		"doesn't respect",
	}
	schemes := []URIScheme{
		fileUriScheme,
		httpUriScheme,
		httpsUriScheme,
	}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response body %s a specific JSON schema", ComponentType, verb),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		for _, opt := range opts {
			for _, scheme := range schemes {
				var phrases []string
				format := fmt.Sprintf(`body %s json schema %s://(.*)`, scheme, verb)
				if opt == opts[0] {
					phrases = createResponseLinePart(format)
				} else {
					phrases = createAliasedResponseLinePart(format)
				}
				f := func(c api.ScenarioContext) any {
					return instance.createThenBodyCompliesWithJsonSchema(c, scheme, FactoryOpts[any]{
						ResolveValueBeforeAssertion: true,
						AssertAlias:                 opt == opts[1],
						AssertTrue:                  verb == verbs[0],
					})
				}
				for _, phrase := range phrases {
					step.Options = append(step.Options, api.Option{
						HandlerFactory: f,
						Regexp:         phrase,
						Description:    step.Description,
					})
				}
			}
		}
		ctx.Then(step)
	}
}

func (instance *ThenHttpResponseStepFactory) createThenBodyCompliesWithJsonSchema(c api.ScenarioContext, scheme URIScheme, opts FactoryOpts[any]) any {
	f := func(alias, value string) error {
		res, err := instance.getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		if res.headers["content-type"] != "application/json" && res.headers["content-type"] != "application/problem+json" {
			return fmt.Errorf("headers[content-type] must be application/json or application/problem+json")
		}
		if scheme == noneUriScheme {
			return fmt.Errorf("URI scheme must be defined")
		}
		if scheme != fileUriScheme && scheme != httpUriScheme && scheme != httpsUriScheme {
			return fmt.Errorf(`URI scheme "%s" isn't supported`, scheme)
		}
		v, err := c.GetValue(ComponentType, "schemes")
		if err != nil {
			return err
		}
		var schemes map[string]gojsonschema.JSONLoader
		if v == nil {
			schemes = make(map[string]gojsonschema.JSONLoader)
			_ = c.SetValue(ComponentType, "schemes", schemes)
		} else {
			schemes = v.(map[string]gojsonschema.JSONLoader)
		}
		valueOf, hasValue := schemes[value]
		if !hasValue {
			binary, prob := instance.doGetFromURI(c, scheme, value)
			if prob != nil {
				return prob
			}
			l := gojsonschema.NewBytesLoader(binary)
			schemes[value] = l
			valueOf = l
		}
		schemaLoader := valueOf.(gojsonschema.JSONLoader)
		documentLoader := gojsonschema.NewBytesLoader(res.body)
		r, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return err
		}
		if opts.AssertTrue {
			if r.Valid() {
				return nil
			}
			return errors.New(
				strings.Join(funk.Map(r.Errors(), func(desc gojsonschema.ResultError) string {
					return fmt.Sprintf("- %s", desc)
				}).([]string), "\n"),
			)
		}
		if !r.Valid() {
			return nil
		}
		return fmt.Errorf("schema respects the schema %s://%s when it shouldn't", scheme, value)
	}
	if !opts.AssertAlias {
		return func(value string) error {
			return f("", value)
		}
	}
	return f
}

func (instance *ThenHttpResponseStepFactory) doGetFromURI(c api.ScenarioContext, scheme URIScheme, value string) ([]byte, error) {
	if scheme == fileUriScheme {
		return c.GetFS().ReadFile(value)
	} else {
		url := string(scheme) + "://" + value
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		h, err := c.GetValue(ComponentType, "httpClient")
		if err != nil {
			return nil, err
		}
		if h == nil {
			h = http.DefaultClient
		}
		res, err := h.(pkgHttp.Client).Do(req)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(res.Body)
	}
}

func (instance *ThenHttpResponseStepFactory) getHttpResponse(c api.ScenarioContext, alias string) (*Response, error) {
	component, err := c.GetGivenComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	req, isHttpRequest := component.(*Request)
	if !isHttpRequest {
		return nil, fmt.Errorf("cannot retrieve %s from context", ComponentType)
	}
	if req.response == nil {
		return nil, fmt.Errorf("%s must be submitted before using the response", ComponentType)
	}
	return req.response, nil
}

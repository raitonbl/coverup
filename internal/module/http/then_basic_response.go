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

const ContentTypeHeaderName = "content-type"

type ThenBasicHttpResponseStepFactory struct {
}

func (instance *ThenBasicHttpResponseStepFactory) New(ctx api.StepDefinitionContext) {
	instance.enableHeadersStepSupport(ctx)
	instance.enableStatusCodeStepSupport(ctx)
	instance.enableResponseBodyStepSupport(ctx)
	instance.enableHeaderComparisonStepSupport(ctx)
}

func (instance *ThenBasicHttpResponseStepFactory) enableHeaderComparisonStepSupport(ctx api.StepDefinitionContext) {
	ops := PathOperations{
		ExpressionPattern:          `(.*)`,
		Line:                       "header",
		ConvertToNumberIfNecessary: true,
		PhraseFactory:              createResponseLinePart,
		AliasedPhraseFactory:       createAliasedResponseLinePart,
		ExtractFromResponse: func(res *Response, expr string) (any, error) {
			return res.headers[expr], nil
		},
	}
	ops.New(ctx)
}

func (instance *ThenBasicHttpResponseStepFactory) enableStatusCodeStepSupport(ctx api.StepDefinitionContext) {
	verbs := []string{"should be", "shouldn't be"}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response status code %s equal to a specific HTTP status code", ComponentType, verb),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		phrases := getResponseLineRegexp(fmt.Sprintf(`status code %s (\d+)$`, verb), opts)

		for _, phrase := range phrases {
			step.Options = append(step.Options, api.Option{
				Regexp:         phrase,
				Description:    step.Description,
				HandlerFactory: instance.statusCodeAssertionFactory(verb == verbs[0], opts[1] == ""),
			})
		}
		ctx.Then(step)
	}
}

func (instance *ThenBasicHttpResponseStepFactory) statusCodeAssertionFactory(assertTrue, assertAlias bool) func(api.ScenarioContext) any {
	return func(c api.ScenarioContext) any {
		return instance.createStatusCodeAssertionHandler(c, FactoryOpts[any]{
			ResolveValueBeforeAssertion: true,
			AssertAlias:                 assertAlias,
			AssertTrue:                  assertTrue,
		})
	}
}

func (instance *ThenBasicHttpResponseStepFactory) createStatusCodeAssertionHandler(c api.ScenarioContext, opts FactoryOpts[any]) any {
	f := func(alias string, statusCode float64) error {
		req, err := getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		if opts.AssertTrue && req.statusCode == statusCode || !opts.AssertTrue && req.statusCode != statusCode {
			return nil
		}
		return fmt.Errorf("status code %v, but got %v", statusCode, req.statusCode)
	}
	if !opts.AssertAlias {
		return func(statusCode float64) error {
			return f("", statusCode)
		}
	}
	return f
}

func (instance *ThenBasicHttpResponseStepFactory) enableHeadersStepSupport(ctx api.StepDefinitionContext) {
	instance.thenHeadersContains(ctx)
	instance.thenHeadersEqualsTo(ctx)
}

func (instance *ThenBasicHttpResponseStepFactory) thenHeadersContains(ctx api.StepDefinitionContext) {
	instance.defineHeaderAssertions(ctx, instance.headersContainsPredicate, [][]string{
		{"should contain", "has exact headers"},
		{"shouldn't contain", "doesn't have exact headers"},
	})
}

func (instance *ThenBasicHttpResponseStepFactory) thenHeadersEqualsTo(ctx api.StepDefinitionContext) {
	instance.defineHeaderAssertions(ctx, instance.headersEqualPredicate, [][]string{
		{"should be", "has exact headers"},
		{"shouldn't be", "doesn't have exact headers"},
	})
}

func (instance *ThenBasicHttpResponseStepFactory) defineHeaderAssertions(ctx api.StepDefinitionContext, predicate func(map[string]string, map[string]string) bool, verbs [][]string) {
	for _, entry := range verbs {
		verb := entry[0]
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a specific %s response %s", ComponentType, entry[1]),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		phrases := getResponseLineRegexp(fmt.Sprintf(`headers %s:$`, verb), opts)

		for _, phrase := range phrases {
			step.Options = append(step.Options, api.Option{
				Regexp:         phrase,
				Description:    step.Description,
				HandlerFactory: instance.headersAssertionFactory(predicate, verb == verbs[0][0], opts[1] == ""),
			})
		}
		ctx.Then(step)
	}
}

func (instance *ThenBasicHttpResponseStepFactory) headersAssertionFactory(predicate func(map[string]string, map[string]string) bool, assertTrue, assertAlias bool) func(api.ScenarioContext) any {
	return func(c api.ScenarioContext) any {
		return instance.createHeadersAssertionHandler(c, FactoryOpts[any]{
			ResolveValueBeforeAssertion: true,
			AssertAlias:                 assertAlias,
			AssertTrue:                  assertTrue,
		}, predicate)
	}
}

func (instance *ThenBasicHttpResponseStepFactory) createHeadersAssertionHandler(c api.ScenarioContext, opts FactoryOpts[any], predicate func(map[string]string, map[string]string) bool) any {
	f := func(alias string, table *godog.Table) error {
		res, err := getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		definition := make(map[string]string)
		for _, row := range table.Rows {
			k := row.Cells[0].Value
			v := row.Cells[1].Value
			if opts.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(v)
				if prob != nil {
					return prob
				}
				definition[k] = fmt.Sprintf("%v", valueOf)
				continue
			}
			definition[k] = v
		}
		if opts.AssertTrue == predicate(res.headers, definition) {
			return nil
		}
		return instance.headersMismatchError(res.headers)
	}
	if !opts.AssertAlias {
		return func(table *godog.Table) error {
			return f("", table)
		}
	}
	return f
}

func (instance *ThenBasicHttpResponseStepFactory) headersContainsPredicate(fromResponse, definition map[string]string) bool {
	for k, v := range definition {
		if fromResponse[k] != v {
			return false
		}
	}
	return true
}

func (instance *ThenBasicHttpResponseStepFactory) headersEqualPredicate(fromResponse, definition map[string]string) bool {
	if len(fromResponse) != len(definition) {
		return false
	}
	for k, v := range fromResponse {
		if definition[k] != v {
			return false
		}
	}
	return true
}

func (instance *ThenBasicHttpResponseStepFactory) headersMismatchError(headers map[string]string) error {
	sb := strings.Builder{}
	for k, v := range headers {
		sb.WriteString(k + "=" + v + "\n")
	}
	return fmt.Errorf("response headers:\n%s", sb.String())
}

func (instance *ThenBasicHttpResponseStepFactory) enableResponseBodyStepSupport(ctx api.StepDefinitionContext) {
	instance.thenBodyEqualTo(ctx)
	instance.thenBodyEqualFile(ctx)
	instance.thenBodyRespectsJsonSchema(ctx)
}

func (instance *ThenBasicHttpResponseStepFactory) thenBodyEqualTo(ctx api.StepDefinitionContext) {
	instance.thenBodyEqualToSrc(ctx, "body", ":", false)
}

func (instance *ThenBasicHttpResponseStepFactory) thenBodyEqualFile(ctx api.StepDefinitionContext) {
	instance.thenBodyEqualToSrc(ctx, "file", " file://(.*)", true)
}

func (instance *ThenBasicHttpResponseStepFactory) thenBodyEqualToSrc(ctx api.StepDefinitionContext, targetType, regex string, isFileURI bool) {
	verbs := []string{"should be", "shouldn't be"}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response body %s equal to a specific %s", ComponentType, verb, targetType),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		phrases := getResponseLineRegexp(fmt.Sprintf(`body %s%s`, verb, regex), opts)

		for _, phrase := range phrases {
			step.Options = append(step.Options, api.Option{
				HandlerFactory: instance.bodyEqualsToAssertionFactory(isFileURI, verb == verbs[0], opts[1] == ""),
				Regexp:         phrase,
				Description:    step.Description,
			})
		}
		ctx.Then(step)
	}
}

func (instance *ThenBasicHttpResponseStepFactory) bodyEqualsToAssertionFactory(isFileURI, assertTrue, assertAlias bool) func(api.ScenarioContext) any {
	return func(c api.ScenarioContext) any {
		return instance.createBodyEqualsToAssertionHandler(c, isFileURI, FactoryOpts[any]{
			ResolveValueBeforeAssertion: true,
			AssertAlias:                 assertAlias,
			AssertTrue:                  assertTrue,
		})
	}
}

func (instance *ThenBasicHttpResponseStepFactory) createBodyEqualsToAssertionHandler(c api.ScenarioContext, isFileURI bool, opts FactoryOpts[any]) any {
	if isFileURI {
		f := func(alias, value string) error {
			binary, err := c.GetFS().ReadFile(value)
			if err != nil {
				return err
			}
			return instance.assertResponseBodyEqualsTo(c, alias, binary, opts)
		}
		if !opts.AssertAlias {
			return func(value string) error {
				return f("", value)
			}
		}
		return f
	}
	f := func(alias string, value *godog.DocString) error {
		return instance.assertResponseBodyEqualsTo(c, alias, []byte(value.Content), opts)
	}
	if !opts.AssertAlias {
		return func(value *godog.DocString) error {
			return f("", value)
		}
	}
	return f
}

func (instance *ThenBasicHttpResponseStepFactory) assertResponseBodyEqualsTo(c api.ScenarioContext, alias string, binary []byte, opts FactoryOpts[any]) error {
	res, err := getHttpResponse(c, alias)
	if err != nil {
		return err
	}
	var predicate func() (bool, error)
	if res.headers[ContentTypeHeaderName] == "application/json" {
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
	return fmt.Errorf("response shouldn't match expectation")
}

func (instance *ThenBasicHttpResponseStepFactory) thenBodyRespectsJsonSchema(ctx api.StepDefinitionContext) {
	verbs := []string{"respects", "doesn't respect"}
	schemes := []URIScheme{fileUriScheme, httpUriScheme, httpsUriScheme}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a %s response body %s a specific JSON schema", ComponentType, verb),
			Options:     nil,
		}
		opts := []string{"", httpRequestRegex}
		for _, scheme := range schemes {
			phrases := getResponseLineRegexp(fmt.Sprintf(`body %s json schema %s://(.*)`, verb, scheme), opts)

			for _, phrase := range phrases {
				step.Options = append(step.Options, api.Option{
					HandlerFactory: instance.jsonSchemaAssertionFactory(scheme, verb == verbs[0], opts[1] == ""),
					Regexp:         phrase,
					Description:    step.Description,
				})
			}
		}
		ctx.Then(step)
	}
}

func (instance *ThenBasicHttpResponseStepFactory) jsonSchemaAssertionFactory(scheme URIScheme, assertTrue, assertAlias bool) func(api.ScenarioContext) any {
	return func(c api.ScenarioContext) any {
		return instance.createThenBodyCompliesWithJsonSchema(c, scheme, FactoryOpts[any]{
			ResolveValueBeforeAssertion: true,
			AssertAlias:                 assertAlias,
			AssertTrue:                  assertTrue,
		})
	}
}

func (instance *ThenBasicHttpResponseStepFactory) createThenBodyCompliesWithJsonSchema(c api.ScenarioContext, scheme URIScheme, opts FactoryOpts[any]) any {
	f := func(alias, value string) error {
		res, err := getHttpResponse(c, alias)
		if err != nil {
			return err
		}
		if res.headers[ContentTypeHeaderName] != "application/json" && res.headers[ContentTypeHeaderName] != "application/problem+json" {
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
			binary, prob := instance.fetchContentFromURI(c, scheme, value)
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

func (instance *ThenBasicHttpResponseStepFactory) fetchContentFromURI(c api.ScenarioContext, scheme URIScheme, value string) ([]byte, error) {
	if scheme == fileUriScheme {
		return c.GetFS().ReadFile(value)
	}
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

func getHttpResponse(c api.ScenarioContext, alias string) (*Response, error) {
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

func getResponseLineRegexp(format string, opts []string) []string {
	var phrases []string
	for _, opt := range opts {
		if opt == opts[0] {
			phrases = append(phrases, createResponseLinePart(format)...)
		} else {
			phrases = append(phrases, createAliasedResponseLinePart(format)...)
		}
	}
	return phrases
}

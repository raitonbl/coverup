package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	pkgHttp "github.com/raitonbl/coverup/pkg/http"
	"github.com/thoas/go-funk"
	"io"
	"net/http"
	"strings"
)

const (
	expressionShouldBeStringErrorf = "%s should be a string but got %v"
)

var methods = []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "POST", "PATCH", "DELETE"}

type GivenHttpRequestStepFactory struct {
}

func (instance *GivenHttpRequestStepFactory) New(ctx api.StepDefinitionContext) {
	instance.given(ctx)
	// Configure Request
	instance.givenHeaders(ctx)
	instance.givenMethod(ctx)
	instance.givenHeader(ctx)
	instance.givenPath(ctx)
	instance.givenURL(ctx)
	// Request Body
	instance.givenRequestBody(ctx)
	instance.givenRequestForm(ctx)
	// Submit Request
	instance.thenSubmitRequest(ctx)
}

func (instance *GivenHttpRequestStepFactory) thenSubmitRequest(ctx api.StepDefinitionContext) {
	phrases := []string{`client submits the %s`, `%s submits the %s`}
	args := []string{
		ComponentType,
		httpRequestRegex,
	}
	for _, phrase := range phrases {
		step := api.StepDefinition{
			Options:     make([]api.Option, 0),
			Description: fmt.Sprintf("Submits a previously defined %s and stores the response", ComponentType),
		}
		for _, arg := range args {
			description := fmt.Sprintf("Submits the current %s", ComponentType)
			if arg == args[1] {
				description = fmt.Sprintf("Submits the %s whose given name matches the specified", ComponentType)
			}
			var variations []string
			if phrase == phrases[1] {
				description = "On behalf of the specified Entity, " + strings.ToUpper(string(description[0])) + description[1:]
				variations = []string{fmt.Sprintf("^"+phrase, arg), fmt.Sprintf("^(?i)the "+phrase, arg)}
			} else {
				variations = []string{fmt.Sprintf("^(?i)"+phrase, arg), fmt.Sprintf("^(?i)the "+phrase, arg)}
			}
			f := func(c api.ScenarioContext) any {
				if phrase == phrases[0] && arg == args[0] {
					return func() error {
						return instance.doSubmitHttpRequest(c, "", "")
					}
				} else if phrase == phrases[0] && arg == args[1] {
					return func(alias string) error {
						return instance.doSubmitHttpRequest(c, "", alias)
					}
				} else if phrase == phrases[1] && arg == args[0] {
					return func(entityId string) error {
						return instance.doSubmitHttpRequest(c, entityId, "")
					}
				} else {
					return func(entityId, alias string) error {
						return instance.doSubmitHttpRequest(c, entityId, alias)
					}
				}
			}
			for _, variation := range variations {
				step.Options = append(step.Options, api.Option{
					HandlerFactory: f,
					Regexp:         variation,
					Description:    description,
				})
			}
		}
		ctx.When(step)
	}

}

func (instance *GivenHttpRequestStepFactory) doSubmitHttpRequest(c api.ScenarioContext, onBehalfOf string, alias string) error {
	req, err := instance.getHttpRequest(c, alias)
	if err != nil {
		return err
	}
	if onBehalfOf != "" {
		if err = instance.setOnBehalfOf(c, req, onBehalfOf); err != nil {
			return err
		}
	}
	var body io.Reader
	if req.body != nil {
		body = bytes.NewReader(req.body)
	}
	serverURI := req.serverURL
	if req.path != "" {
		serverURI += req.path
	}
	httpClientRequest, err := http.NewRequest(req.method, serverURI, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	for k, v := range req.headers {
		httpClientRequest.Header.Set(k, v)
	}
	valueOf, err := c.GetValue(ComponentType, "httpClient")
	if err != nil {
		return err
	}
	httpClient, isHttpClient := valueOf.(pkgHttp.Client)
	if !isHttpClient || httpClient == nil {
		httpClient = http.DefaultClient
		_ = c.SetValue(ComponentType, "httpClient", httpClient)
	}
	res, err := httpClient.Do(httpClientRequest)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	binary, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	headers := make(map[string]string)
	for k, v := range res.Header {
		headers[k] = strings.Join(v, ",")
	}
	req.response = &Response{
		body:       binary,
		headers:    headers,
		statusCode: float64(res.StatusCode),
		pathCache:  make(map[string]any),
	}
	return nil
}

func (instance *GivenHttpRequestStepFactory) givenRequestForm(ctx api.StepDefinitionContext) {
	instance.withFormFile(ctx)
	instance.withFormEncType(ctx)
	instance.withFormAttribute(ctx)
	instance.withFormAttributes(ctx)
}

func (instance *GivenHttpRequestStepFactory) withFormEncType(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies the HTTP Form encription type for the current %s", ComponentType),
	}
	f := func(c api.ScenarioContext) any {
		return func(value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			if value == "multipart/form-data" || value == "application/x-www-form-urlencoded" {
				req.body = nil
				if req.form == nil {
					req.form = &Form{}
				}
				req.form.encType = value
				return nil
			}
			return fmt.Errorf(`encType "%s"" is not supported for form`, value)
		}
	}
	for _, regexp := range createRequestLinePart(`http request form enctype is ([^"]*)$`) {
		step.Options = append(step.Options, api.Option{
			Regexp:         regexp,
			Description:    step.Description,
			HandlerFactory: f,
		})
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) withFormAttribute(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies a HTTP Form attribute for the current %s", ComponentType),
	}
	f := func(c api.ScenarioContext) any {
		return func(name, value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			if err = instance.doGivenFormAttribute(c, req, name, value, false); err != nil {
				return err
			}
			return nil
		}
	}
	args := []string{
		`"([^"]+)"`,
		fmt.Sprintf(`"%s"`, api.ValueExpression),
	}
	for _, arg := range args {
		for _, regexp := range createRequestLinePart(fmt.Sprintf(`http request form attribute "([a-zA-Z_]+)" is %s$`, arg)) {
			step.Options = append(step.Options, api.Option{
				Regexp:         regexp,
				Description:    step.Description,
				HandlerFactory: f,
			})
		}
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) withFormFile(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies a HTTP Form attribute, the specified file , for the current %s", ComponentType),
	}
	f := func(c api.ScenarioContext) any {
		return func(name, value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			if err = instance.doGivenFormAttribute(c, req, name, "file://"+value, true); err != nil {
				return err
			}
			return nil
		}
	}
	args := []string{
		`([^"]+)`,
		fmt.Sprintf(`%s`, api.ValueExpression),
	}
	for _, arg := range args {
		for _, regexp := range createRequestLinePart(fmt.Sprintf(`http request form attribute "([a-zA-Z_]+)" is file://%s$`, arg)) {
			step.Options = append(step.Options, api.Option{
				Regexp:         regexp,
				Description:    step.Description,
				HandlerFactory: f,
			})
		}
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) withFormAttributes(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies the HTTP Form attributes for the current %s", ComponentType),
	}
	f := func(c api.ScenarioContext) any {
		return func(table *godog.Table) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			req.body = nil
			if req.form == nil {
				req.form = &Form{}
			}
			req.form.attributes = make(map[string]any)
			for _, row := range table.Rows {
				if err = instance.doGivenFormAttribute(c, req, row.Cells[0].Value, row.Cells[1].Value, true); err != nil {
					return err
				}
			}
			return nil
		}
	}
	for _, regexp := range createRequestLinePart(`http request form attributes are:$`) {
		step.Options = append(step.Options, api.Option{
			Regexp:         regexp,
			Description:    step.Description,
			HandlerFactory: f,
		})
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) givenRequestBody(ctx api.StepDefinitionContext) {
	stepVariations := []string{
		`:$`,
		` file://(.+)$`,
	}
	for _, variation := range stepVariations {
		var extract func(api.ScenarioContext, any) ([]byte, error)
		description := ""
		if variation == stepVariations[0] {
			extract = func(_ api.ScenarioContext, value any) ([]byte, error) {
				return []byte(value.(*godog.DocString).Content), nil
			}
			description = fmt.Sprintf("Specifies the HTTP Request body for the current %s", ComponentType)
		} else {
			extract = func(c api.ScenarioContext, value any) ([]byte, error) {
				return c.GetFS().ReadFile(value.(string))
			}
			description = fmt.Sprintf("Specifies the HTTP Request body, reading from the specified file, for the current %s", ComponentType)
		}
		step := api.StepDefinition{
			Description: description,
			Options:     make([]api.Option, 0),
		}
		args := []string{serverURLRegex, api.PropertyExpression}
		f := func(c api.ScenarioContext) any {
			doExec := func(value any) error {
				req, err := instance.getHttpRequest(c, "")
				if err != nil {
					return err
				}
				var v any
				if variation != stepVariations[0] {
					valueOf, prob := c.Resolve(value.(string))
					if prob != nil {
						return prob
					}
					x, isString := valueOf.(string)
					if !isString {
						return fmt.Errorf(expressionShouldBeStringErrorf, value, valueOf)
					}
					v = x
				}

				binary, err := extract(c, v)
				if err != nil {
					return err
				}
				req.body = binary
				return nil
			}
			if variation == stepVariations[0] {
				return func(value *godog.DocString) error {
					return doExec(value)
				}
			}
			return func(value string) error {
				return doExec(value)
			}
		}
		for _, prefix := range createRequestLinePart(fmt.Sprintf(`http request body%s`, variation)) {
			for _, parameter := range args {
				step.Options = append(step.Options, api.Option{
					Regexp:         fmt.Sprintf(`%s %s`, prefix, parameter),
					Description:    step.Description,
					HandlerFactory: f,
				})
			}
		}
		ctx.Given(step)
	}
}

func (instance *GivenHttpRequestStepFactory) givenURL(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies the HTTP Server URL for the current %s", ComponentType),
	}
	args := []string{serverURLRegex, api.PropertyExpression}
	f := func(c api.ScenarioContext) any {
		return func(value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			valueOf, err := c.Resolve(value)
			if err != nil {
				return err
			}
			v, isString := valueOf.(string)
			if !isString {
				return fmt.Errorf(expressionShouldBeStringErrorf, value, valueOf)
			}
			req.serverURL = v
			return nil
		}
	}
	for _, prefix := range createRequestLinePart(`http request URL is`) {
		for _, parameter := range args {
			step.Options = append(step.Options, api.Option{
				Regexp:         fmt.Sprintf(`%s %s`, prefix, parameter),
				Description:    step.Description,
				HandlerFactory: f,
			})
		}
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) givenPath(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Description: fmt.Sprintf("Specifies the HTTP path for the current %s", ComponentType),
		Options:     make([]api.Option, 0),
	}
	args := []string{relativeURIRegex, api.ValueExpression}
	f := func(c api.ScenarioContext) any {
		return func(value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			valueOf, err := c.Resolve(value)
			if err != nil {
				return err
			}
			v, isString := valueOf.(string)
			if !isString {
				return fmt.Errorf(expressionShouldBeStringErrorf, value, valueOf)
			}
			req.path = "/" + v
			return nil
		}
	}
	for _, prefix := range createRequestLinePart(`http request path is`) {
		for _, parameter := range args {
			step.Options = append(step.Options, api.Option{
				Regexp:         fmt.Sprintf(`%s %s`, prefix, parameter),
				Description:    step.Description,
				HandlerFactory: f,
			})
		}
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) given(ctx api.StepDefinitionContext) {
	ctx.Given(api.StepDefinition{
		Description: "Defines a http request into the Scenario",
		Options: []api.Option{
			{
				HandlerFactory: instance.createHttpRequest(),
				Regexp:         `^(?i)a HttpRequest$`,
				Description:    "Defines a http request into the scenario so it can be configured",
			},
			{
				HandlerFactory: instance.createNamedHttpRequest(),
				Regexp:         `^(?i)a HttpRequest named (.+)$`,
				Description:    "Defines a http request, with given name, into the scenario so it can be configured",
			},
		},
	})
}

func (instance *GivenHttpRequestStepFactory) givenHeader(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Options:     make([]api.Option, 0),
		Description: fmt.Sprintf("Specifies an HTTP header for the current %s", ComponentType),
	}
	f := func(c api.ScenarioContext) any {
		return func(value string) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			valueOf, err := c.Resolve(value)
			if err != nil {
				return err
			}
			v, isString := valueOf.(string)
			if !isString {
				return fmt.Errorf(expressionShouldBeStringErrorf, value, valueOf)
			}
			req.path = v
			return nil
		}
	}
	for _, regex := range createRequestLinePart(`http request header (.*) is "([^"]*)"$`) {
		step.Options = append(step.Options, api.Option{
			Regexp:         regex,
			Description:    step.Description,
			HandlerFactory: f,
		})
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) createHttpRequest() api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		return func() error {
			return instance.doGivenHttpRequest(c, "")
		}
	}
}

func (instance *GivenHttpRequestStepFactory) createNamedHttpRequest() api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		return func(alias string) error {
			return instance.doGivenHttpRequest(c, alias)
		}
	}
}

func (instance *GivenHttpRequestStepFactory) doGivenHttpRequest(c api.ScenarioContext, alias string) error {
	return c.AddGivenComponent(ComponentType, &Request{
		headers: make(map[string]string),
	}, alias)
}

func (instance *GivenHttpRequestStepFactory) givenMethod(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Description: fmt.Sprintf("Specifies the HTTP method for the current %s", ComponentType),
		Options:     make([]api.Option, 0),
	}
	f := func(c api.ScenarioContext) any {
		return func(method string) error {
			if !funk.Contains(methods, method) {
				return fmt.Errorf("cannot assign %s.method to %s", ComponentType, method)
			}
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			req.method = method
			return nil
		}
	}
	for _, regex := range createRequestLinePart(`http request method is ([^"]*)$`) {
		step.Options = append(step.Options, api.Option{
			Regexp:         regex,
			Description:    step.Description,
			HandlerFactory: f,
		})
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) givenHeaders(ctx api.StepDefinitionContext) {
	step := api.StepDefinition{
		Description: fmt.Sprintf("Specifies the HTTP headers for the current %s", ComponentType),
		Options:     make([]api.Option, 0),
	}
	f := func(c api.ScenarioContext) any {
		return func(table *godog.Table) error {
			req, err := instance.getHttpRequest(c, "")
			if err != nil {
				return err
			}
			if req.headers == nil {
				req.headers = make(map[string]string)
			}
			req.headers = make(map[string]string)
			for _, row := range table.Rows {
				if err = instance.doGivenHeader(c, req, row.Cells[0].Value, row.Cells[1].Value); err != nil {
					return err
				}
			}
			return nil
		}
	}
	for _, regex := range createRequestLinePart(`http request headers are:$`) {
		step.Options = append(step.Options, api.Option{
			Regexp:         regex,
			Description:    step.Description,
			HandlerFactory: f,
		})
	}
	ctx.Given(step)
}

func (instance *GivenHttpRequestStepFactory) getHttpRequest(c api.ScenarioContext, alias string) (*Request, error) {
	component, err := c.GetGivenComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	req, isHttpRequest := component.(*Request)
	if !isHttpRequest {
		return nil, fmt.Errorf("cannot retrieve %s from context", ComponentType)
	}
	if req.response != nil {
		return nil, fmt.Errorf("cannot modify a %s that has been submitted", ComponentType)
	}
	return req, nil
}

func (instance *GivenHttpRequestStepFactory) doGivenHeader(c api.ScenarioContext, req *Request, key string, value string) error {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	valueOf, err := c.Resolve(value)
	if err != nil {
		return err
	}
	req.headers[key] = fmt.Sprintf("%v", valueOf)
	return nil
}

func (instance *GivenHttpRequestStepFactory) doGivenFormAttribute(c api.ScenarioContext, req *Request, name, value string, mightBeFileURI bool) error {
	valueOf, err := c.Resolve(value)
	if err != nil {
		return err
	}
	if mightBeFileURI && strings.HasPrefix(value, "file://") {
		binary, prob := c.GetFS().ReadFile(value)
		if prob != nil {
			return prob
		}
		req.form.attributes[name] = binary
		return nil
	}
	v, isString := valueOf.(string)
	if !isString {
		return fmt.Errorf(expressionShouldBeStringErrorf, value, valueOf)
	}
	req.body = nil
	if req.form == nil {
		req.form = &Form{}
	}
	req.form.attributes[name] = fmt.Sprintf("%v", v)
	return nil
}

func (instance *GivenHttpRequestStepFactory) setOnBehalfOf(c api.ScenarioContext, req *Request, onBehalfOf string) error {
	valueOf, err := c.Resolve(fmt.Sprintf(`{{Entities.%s}}`, onBehalfOf))
	if err != nil {
		return err
	}
	if bearerToken, isBearerToken := valueOf.(api.BearerToken); isBearerToken {
		req.headers["Authorization"] = fmt.Sprintf(`Bearer %s`, bearerToken.Value)
		return nil
	}
	if usernameAndPassword, isUsernameAndPassword := valueOf.(api.UsernameAndPassword); isUsernameAndPassword {
		value := fmt.Sprintf(`%s:%s`, usernameAndPassword.Username, usernameAndPassword.Password)
		req.headers["Authorization"] = base64.StdEncoding.EncodeToString([]byte(value))
		return nil
	}
	return fmt.Errorf(`entity %s ins't supported'`, onBehalfOf)
}

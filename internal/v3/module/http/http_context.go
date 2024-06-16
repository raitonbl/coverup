package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/v3/pkg"
	"github.com/thoas/go-funk"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const fileURISchema = "file://"
const ComponentType = "HttpRequest"

var pathRegexp, _ = regexp.Compile(`^\$((\.\w+)|(\[\d+\]))*$`)

type Scenario struct {
	schemas map[string]any
	ctx     pkg.ScenarioContext
}

func (instance *Scenario) WithRequest() error {
	return instance.WithRequestWhenAlias("")
}

func (instance *Scenario) WithRequestWhenAlias(alias string) error {
	return instance.ctx.Register(ComponentType, &Request{
		headers: make(map[string]string),
	}, alias)
}

// Deprecated
func (instance *Scenario) onHttpRequest(f func(*Request) error) error {
	return instance.onNamedHttpRequest("", f)
}

// Deprecated
func (instance *Scenario) onNamedHttpRequest(alias string, f func(*Request) error) error {
	valueOf, err := instance.ctx.GetComponent(ComponentType, alias)
	if err != nil {
		return err
	}
	if r, isHttpRequest := valueOf.(*Request); isHttpRequest {
		return f(r)
	}
	return fmt.Errorf("please define %s before using it", ComponentType)
}

func (instance *Scenario) getRequest(alias string) (*Request, error) {
	valueOf, err := instance.ctx.GetComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	if r, isHttpRequest := valueOf.(*Request); isHttpRequest {
		return r, nil
	}
	return nil, fmt.Errorf("please define %s before using it", ComponentType)
}

func (instance *Scenario) WithHeaders(table *godog.Table) error {
	return instance.onHttpRequest(func(req *Request) error {
		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		for _, row := range table.Rows {
			key := row.Cells[0].Value
			valueOf, prob := instance.ctx.GetValue(row.Cells[1].Value)
			if prob != nil {
				return fmt.Errorf(`cannot determine header "%s" value. caused by:%v"`, key, prob)
			}
			err := instance.withHeader(req, key, fmt.Sprintf("%v", valueOf))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (instance *Scenario) WithHeader(name, value string) error {
	return instance.onHttpRequest(func(req *Request) error {
		return instance.withHeader(req, name, value)
	})
}

func (instance *Scenario) withHeader(req *Request, name, value string) error {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	valueOf, err := instance.ctx.GetValue(value)
	if err != nil {
		return err
	}
	req.headers[name] = fmt.Sprintf("%v", valueOf)
	return nil
}

func (instance *Scenario) WithMethod(method string) error {
	return instance.onHttpRequest(func(req *Request) error {
		switch method {
		case "OPTIONS":
			return instance.withMethod(req, method)
		case "HEAD":
			return instance.withMethod(req, method)
		case "GET":
			return instance.withMethod(req, method)
		case "PUT":
			return instance.withMethod(req, method)
		case "POST":
			return instance.withMethod(req, method)
		case "PATCH":
			return instance.withMethod(req, method)
		case "DELETE":
			return instance.withMethod(req, method)
		default:
			return fmt.Errorf("cannot assign %s.method to %s ", ComponentType, method)
		}
	})
}

func (instance *Scenario) withMethod(req *Request, method string) error {
	req.method = method
	return nil
}

func (instance *Scenario) WithPath(path string) error {
	return instance.onHttpRequest(func(req *Request) error {
		valueOf, err := instance.ctx.GetValue(path)
		if err != nil {
			return err
		}
		req.path = fmt.Sprintf("/%v", valueOf)
		return nil
	})
}

func (instance *Scenario) WithHttpPath(url string) error {
	return instance.onHttpRequest(func(req *Request) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("http://%v", valueOf)
		return nil
	})
}

func (instance *Scenario) WithHttpsPath(url string) error {
	return instance.onHttpRequest(func(req *Request) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("https://%v", valueOf)
		return nil
	})
}

func (instance *Scenario) WithServerURL(url string) error {
	return instance.onHttpRequest(func(req *Request) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("%v", valueOf)
		return nil
	})
}

func (instance *Scenario) WithBody(body *godog.DocString) error {
	return instance.withBody(func() ([]byte, error) {
		return []byte(body.Content), nil
	})
}

func (instance *Scenario) withBody(s func() ([]byte, error)) error {
	return instance.onHttpRequest(func(req *Request) error {
		binary, err := s()
		if err != nil {
			return err
		}
		req.form = nil
		req.body = binary
		return nil
	})
}

func (instance *Scenario) WithBodyFileURI(value string) error {
	return instance.withBody(func() ([]byte, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.ctx.GetFS().ReadFile(valueOf.(string))
	})
}

func (instance *Scenario) WithAcceptHeader(value string) error {
	return instance.WithHeader("Accept", value)
}

func (instance *Scenario) WithContentTypeHeader(value string) error {
	return instance.WithHeader("Content-Type", value)
}

func (instance *Scenario) WithFormEncType(value string) error {
	return instance.onForm(func(form *Form) error {
		if value == "multipart/form-data" || value == "application/x-www-form-urlencoded" {
			form.encType = value
			return nil
		}
		return fmt.Errorf("encType %s is not supported", value)
	})
}

func (instance *Scenario) WithFormAttribute(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v", valueOf), nil
	})
}

func (instance *Scenario) WithFormFile(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.ctx.GetFS().ReadFile(valueOf.(string))
	})
}

func (instance *Scenario) onFormAttribute(name string, f func() (any, error)) error {
	return instance.onForm(func(form *Form) error {
		valueOf, err := f()
		if err != nil {
			return err
		}
		form.attributes[name] = valueOf
		return nil
	})
}

func (instance *Scenario) onForm(s func(*Form) error) error {
	return instance.onHttpRequest(func(req *Request) error {
		form := req.form
		if form == nil {
			req.form = &Form{
				encType:    "multipart/form-data",
				attributes: make(map[string]any),
			}
			form = req.form
		}
		req.body = nil
		return s(form)
	})
}

func (instance *Scenario) get(alias string) (*Request, error) {
	valueOf, err := instance.ctx.GetComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	req, isHttpRequest := valueOf.(*Request)
	if !isHttpRequest {
		if alias == "" {
			return nil, fmt.Errorf(`%s is undefined and needs to be defined before referencing it`, ComponentType)
		}
		return nil, fmt.Errorf(`%s["%s"] is undefined and needs to be defined before referencing it`, ComponentType, alias)
	}
	return req, nil
}

// Deprecated
func (instance *Scenario) SubmitHttpRequest() error {
	return instance.SubmitHttpRequestOnBehalfOfEntity("")
}

// Deprecated
func (instance *Scenario) SubmitHttpRequestOnBehalfOfEntity(id string) error {
	return instance.SubmitNamedHttpRequestOnBehalfOfEntity("", id)
}

// Deprecated
func (instance *Scenario) SubmitNamedHttpRequest(alias string) error {
	return instance.SubmitNamedHttpRequestOnBehalfOfEntity(alias, "")
}

// Deprecated
func (instance *Scenario) SubmitNamedHttpRequestOnBehalfOfEntity(alias, id string) error {
	return instance.onNamedHttpRequest(alias, func(req *Request) error {
		if id != "" {
			panic("not implemented")
		}
		if req.response != nil {
			return fmt.Errorf(`cannot submit the same %s more than once`, ComponentType)
		}
		return instance.doSubmitHttpRequest(req)
	})
}

// Deprecated
func (instance *Scenario) doSubmitHttpRequest(src *Request) error {
	var body io.Reader
	if src.body != nil {
		body = bytes.NewReader(src.body)
	}
	serverURI := src.serverURL
	if src.path != "" {
		serverURI += src.path
	}
	req, err := http.NewRequest(src.method, serverURI, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	for k, v := range src.headers {
		req.Header.Set(k, v)
	}
	httpClient := instance.ctx.GetHttpClient()
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	res, err := httpClient.Do(req)
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
	src.response = &Response{
		body:       binary,
		headers:    headers,
		statusCode: float64(res.StatusCode),
		pathCache:  make(map[string]any),
	}
	return nil
}

// Deprecated
func (instance *Scenario) onNamedResponse(alias string, f func(*Request, *Response) error) error {
	return instance.onNamedHttpRequest(alias, func(req *Request) error {
		if req.response == nil {
			if alias == "" {
				return fmt.Errorf(`%s needs to be submitted before making assertions`, ComponentType)
			} else {
				return fmt.Errorf(`%s["%s"] needs to be submitted before making assertions`, ComponentType, alias)
			}
		}
		return f(req, req.response)
	})
}

func (instance *Scenario) getResponse(alias string) (*Response, error) {
	req, err := instance.getRequest(alias)
	if err != nil {
		return nil, err
	}
	if req.response == nil {
		if alias == "" {
			return nil, fmt.Errorf(`%s needs to be submitted before making assertions`, ComponentType)
		} else {
			return nil, fmt.Errorf(`%s["%s"] needs to be submitted before making assertions`, ComponentType, alias)
		}
	}
	return req.response, nil
}

func (instance *Scenario) AssertResponseStatusCode(statusCode int) error {
	return instance.AssertNamedHttpRequestResponseStatusCode("", statusCode)
}

func (instance *Scenario) AssertNamedHttpRequestResponseStatusCode(alias string, statusCode int) error {
	return instance.onNamedResponse(alias, func(_ *Request, response *Response) error {
		if response.statusCode != float64(statusCode) {
			if alias == "" {
				return fmt.Errorf("%s.StatusCode should be %d but instead got %d", ComponentType, statusCode, response.statusCode)
			}
			return fmt.Errorf(`%s["%s"].StatusCode should be %d but instead got %d`, ComponentType, alias, statusCode, response.statusCode)
		}
		return nil
	})
}

func (instance *Scenario) AssertNamedHttpRequestResponseExactHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *Scenario) AssertResponseExactHeaders(table *godog.Table) error {
	return nil
}

func (instance *Scenario) AssertResponseContainsHeaders(table *godog.Table) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponseContainsHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *Scenario) AssertResponseHeader(name, value string) error {
	return nil
}

func (instance *Scenario) AssertResponseContentType(value string) error {
	return instance.AssertResponseHeader("content-type", value)
}

func (instance *Scenario) AssertNamedHttpRequestResponseContentType(alias, value string) error {
	return instance.AssertNamedHttpRequestResponseHeader(alias, "content-type", value)
}

func (instance *Scenario) AssertNamedHttpRequestResponseHeader(alias string, header string, value string) error {
	return nil
}

func (instance *Scenario) AssertResponseIsValidAgainstSchema(file string) error {
	return instance.AssertNamedHttpRequestResponseIsValidAgainstSchema("", file)
}

// Deprecated
func (instance *Scenario) onNamedHttpRequestResponseWithJsonContentType(alias string, f func(*Request, *Response) error) error {
	return instance.onNamedResponse(alias, func(req *Request, res *Response) error {
		contentType, hasValue := res.headers["content-type"]
		if !hasValue {
			if alias == "" {
				return fmt.Errorf(`%s.headers["content-type"] must be "application/json" or "aplication/problem+json"`, ComponentType)
			}
			return fmt.Errorf(`%s.%s.headers["content-type"] must be "application/json" or "aplication/problem+json"`, ComponentType, alias)
		}
		if contentType != "application/json" && contentType != "application/problem+json" {
			if alias == "" {
				return fmt.Errorf(`%s.headers["content-type"] must be "application/json" or "aplication/problem+json" but content-type is "%s"`, ComponentType, contentType)
			}
			return fmt.Errorf(`%s.%s.headers["content-type"] must be "application/json" or "aplication/problem+json" but content-type is "%s"`, ComponentType, alias, contentType)
		}
		return f(req, res)
	})
}

func (instance *Scenario) getResponseWithJsonContentType(alias string) (*Response, error) {
	res, err := instance.getResponse(alias)
	if err != nil {
		return nil, err
	}
	contentType, hasValue := res.headers["content-type"]
	if !hasValue {
		if alias == "" {
			return nil, fmt.Errorf(`%s.headers["content-type"] must be "application/json" or "aplication/problem+json"`, ComponentType)
		}
		return nil, fmt.Errorf(`%s.%s.headers["content-type"] must be "application/json" or "aplication/problem+json"`, ComponentType, alias)
	}
	if contentType != "application/json" && contentType != "application/problem+json" {
		if alias == "" {
			return nil, fmt.Errorf(`%s.headers["content-type"] must be "application/json" or "aplication/problem+json" but content-type is "%s"`, ComponentType, contentType)
		}
		return nil, fmt.Errorf(`%s.%s.headers["content-type"] must be "application/json" or "aplication/problem+json" but content-type is "%s"`, ComponentType, alias, contentType)
	}
	return res, nil
}

func (instance *Scenario) getResponseValueFromExpression(alias, t string) (any, error) {
	res, err := instance.getResponse(alias)
	if err != nil {
		return nil, err
	}
	expr := t
	if strings.HasPrefix(expr, ".") {
		expr = expr[1:]
	}
	if len(res.body) == 0 {
		if alias == "" {
			return nil, fmt.Errorf(`cannot determine %s.body.$%s since body is undefined`, ComponentType, expr)
		}
		return nil, fmt.Errorf(`cannot determine %s["%s"].body.$%s since body is undefined`, ComponentType, alias, expr)
	}
	valueOf, err := res.JSONPath(expr)
	if err != nil {
		if alias == "" {
			return nil, fmt.Errorf(`cannot determine %s.body.$%s due to error:\n%v`, ComponentType, expr, err)
		}
		return nil, fmt.Errorf(`cannot determine %s["%s"].body.$%s due to error:\n%v`, ComponentType, alias, expr, err)
	}
	return valueOf, nil
}

func (instance *Scenario) AssertNamedHttpRequestResponseIsValidAgainstSchema(alias, value string) error {
	return instance.onNamedHttpRequestResponseWithJsonContentType(alias, func(_ *Request, response *Response) error {
		if instance.schemas == nil {
			instance.schemas = make(map[string]any)
		}
		valueOf, hasValue := instance.schemas[value]
		if !hasValue {
			binary, err := instance.ctx.GetFS().ReadFile(value)
			if err != nil {
				return err
			}
			l := gojsonschema.NewBytesLoader(binary)
			instance.schemas[value] = l
			valueOf = l
		}
		schemaLoader := valueOf.(gojsonschema.JSONLoader)
		documentLoader := gojsonschema.NewBytesLoader(response.body)
		r, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return err
		}
		if !r.Valid() {
			return errors.New(
				strings.Join(funk.Map(r.Errors(), func(desc gojsonschema.ResultError) string {
					return fmt.Sprintf("- %s", desc)
				}).([]string), "\n"),
			)
		}
		return nil
	})
}

func (instance *Scenario) onNamedHttpRequestResponseBodyPath(t, alias string, f func(*Request, *Response, any) error) error {
	return instance.onNamedHttpRequestResponseWithJsonContentType(alias, func(req *Request, res *Response) error {
		expr := t
		if strings.HasPrefix(expr, ".") {
			expr = expr[1:]
		}
		if res.body == nil {
			if alias == "" {
				return fmt.Errorf(`cannot determine %s.body.$%s since body is undefined`, ComponentType, expr)
			}
			return fmt.Errorf(`cannot determine %s["%s"].body.$%s since body is undefined`, ComponentType, alias, expr)
		}
		valueOf, err := res.JSONPath(expr)
		if err != nil {
			if alias == "" {
				return fmt.Errorf(`cannot determine %s.body.$%s due to error:\n%v`, ComponentType, expr, err)
			}
			return fmt.Errorf(`cannot determine %s["%s"].body.$%s due to error:\n%v`, ComponentType, alias, expr, err)
		}
		return f(req, res, valueOf)
	})
}

func (instance *Scenario) AssertResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathLengthIs(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathLengthIs(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsGreaterThan(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsGreaterThanOrEqualTo(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsLesserThan(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsLesserThanOrEqualTo(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsGreaterThanValue(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsGreaterThanOrEqualToValue(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsLesserThanValue(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsLesserThanOrEqualToValue(k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsGreaterThan(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualTo(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsLesserThan(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsLesserThanOrEqualTo(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsGreaterThanValue(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualToValue(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsLesserThanValue(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsLesserThanOrEqualToValue(alias, k string, value float64) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsInStringArray(k string, value string) error {
	return nil
}

func (instance *Scenario) AssertResponsePathIsInNumericArray(k string, value string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsInStringArray(alias, k string, value string) error {
	return nil
}

func (instance *Scenario) AssertNamedHttpRequestResponsePathIsInNumericArray(alias, k string, value string) error {
	return nil
}

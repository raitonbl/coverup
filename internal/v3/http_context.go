package v3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
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

type HttpContext struct {
	schemas map[string]any
	ctx     ScenarioContext
}

func (instance *HttpContext) WithRequest() error {
	return instance.WithRequestWhenAlias("")
}

func (instance *HttpContext) WithRequestWhenAlias(alias string) error {
	return instance.ctx.Register(ComponentType, &HttpRequest{
		headers: make(map[string]string),
	}, alias)
}

func (instance *HttpContext) onHttpRequest(f func(*HttpRequest) error) error {
	return instance.onNamedHttpRequest("", f)
}

func (instance *HttpContext) onNamedHttpRequest(alias string, f func(*HttpRequest) error) error {
	valueOf, err := instance.ctx.GetComponent(ComponentType, alias)
	if err != nil {
		return err
	}
	if r, isHttpRequest := valueOf.(*HttpRequest); isHttpRequest {
		return f(r)
	}
	return fmt.Errorf("please define %s before using it", ComponentType)
}

func (instance *HttpContext) WithHeaders(table *godog.Table) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
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

func (instance *HttpContext) WithHeader(name, value string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		return instance.withHeader(req, name, value)
	})
}

func (instance *HttpContext) withHeader(req *HttpRequest, name, value string) error {
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

func (instance *HttpContext) WithMethod(method string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
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

func (instance *HttpContext) withMethod(req *HttpRequest, method string) error {
	req.method = method
	return nil
}

func (instance *HttpContext) WithPath(path string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(path)
		if err != nil {
			return err
		}
		req.path = fmt.Sprintf("/%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithHttpPath(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("http://%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithHttpsPath(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("https://%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithServerURL(url string) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		valueOf, err := instance.ctx.GetValue(url)
		if err != nil {
			return err
		}
		req.path = ""
		req.serverURL = fmt.Sprintf("%v", valueOf)
		return nil
	})
}

func (instance *HttpContext) WithBody(body *godog.DocString) error {
	return instance.withBody(func() ([]byte, error) {
		return []byte(body.Content), nil
	})
}

func (instance *HttpContext) withBody(s func() ([]byte, error)) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
		binary, err := s()
		if err != nil {
			return err
		}
		req.form = nil
		req.body = binary
		return nil
	})
}

func (instance *HttpContext) WithBodyFileURI(value string) error {
	return instance.withBody(func() ([]byte, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.ctx.GetFS().ReadFile(valueOf.(string))
	})
}

func (instance *HttpContext) WithAcceptHeader(value string) error {
	return instance.WithHeader("Accept", value)
}

func (instance *HttpContext) WithContentTypeHeader(value string) error {
	return instance.WithHeader("Content-Type", value)
}

func (instance *HttpContext) WithFormEncType(value string) error {
	return instance.onForm(func(form *Form) error {
		if value == "multipart/form-data" || value == "application/x-www-form-urlencoded" {
			form.encType = value
			return nil
		}
		return fmt.Errorf("encType %s is not supported", value)
	})
}

func (instance *HttpContext) WithFormAttribute(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v", valueOf), nil
	})
}

func (instance *HttpContext) WithFormFile(name, value string) error {
	return instance.onFormAttribute(name, func() (any, error) {
		valueOf, err := instance.ctx.GetValue(value)
		if err != nil {
			return nil, err
		}
		return instance.ctx.GetFS().ReadFile(valueOf.(string))
	})
}

func (instance *HttpContext) onFormAttribute(name string, f func() (any, error)) error {
	return instance.onForm(func(form *Form) error {
		valueOf, err := f()
		if err != nil {
			return err
		}
		form.attributes[name] = valueOf
		return nil
	})
}

func (instance *HttpContext) onForm(s func(*Form) error) error {
	return instance.onHttpRequest(func(req *HttpRequest) error {
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

func (instance *HttpContext) get(alias string) (*HttpRequest, error) {
	valueOf, err := instance.ctx.GetComponent(ComponentType, alias)
	if err != nil {
		return nil, err
	}
	req, isHttpRequest := valueOf.(*HttpRequest)
	if !isHttpRequest {
		if alias == "" {
			return nil, fmt.Errorf(`%s is undefined and needs to be defined before referencing it`, ComponentType)
		}
		return nil, fmt.Errorf(`%s["%s"] is undefined and needs to be defined before referencing it`, ComponentType, alias)
	}
	return req, nil
}

func (instance *HttpContext) SubmitHttpRequest() error {
	return instance.SubmitHttpRequestOnBehalfOfEntity("")
}

func (instance *HttpContext) SubmitHttpRequestOnBehalfOfEntity(id string) error {
	return instance.SubmitNamedHttpRequestOnBehalfOfEntity("", id)
}

func (instance *HttpContext) SubmitNamedHttpRequest(alias string) error {
	return instance.SubmitNamedHttpRequestOnBehalfOfEntity(alias, "")
}

func (instance *HttpContext) SubmitNamedHttpRequestOnBehalfOfEntity(alias, id string) error {
	return instance.onNamedHttpRequest(alias, func(req *HttpRequest) error {
		if id != "" {
			panic("not implemented")
		}
		if req.response != nil {
			return fmt.Errorf(`cannot submit the same %s more than once`, ComponentType)
		}
		return instance.doSubmitHttpRequest(req)
	})
}

func (instance *HttpContext) doSubmitHttpRequest(src *HttpRequest) error {
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
	src.response = &HttpResponse{
		body:       binary,
		headers:    headers,
		statusCode: res.StatusCode,
		pathCache:  make(map[string]any),
	}
	return nil
}

func (instance *HttpContext) onNamedResponse(alias string, f func(*HttpRequest, *HttpResponse) error) error {
	return instance.onNamedHttpRequest(alias, func(req *HttpRequest) error {
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

func (instance *HttpContext) AssertResponseStatusCode(statusCode int) error {
	return instance.AssertNamedHttpRequestResponseStatusCode("", statusCode)
}

func (instance *HttpContext) AssertNamedHttpRequestResponseStatusCode(alias string, statusCode int) error {
	return instance.onNamedResponse(alias, func(_ *HttpRequest, response *HttpResponse) error {
		if response.statusCode != statusCode {
			if alias == "" {
				return fmt.Errorf("%s.StatusCode should be %d but instead got %d", ComponentType, statusCode, response.statusCode)
			}
			return fmt.Errorf(`%s["%s"].StatusCode should be %d but instead got %d`, ComponentType, alias, statusCode, response.statusCode)
		}
		return nil
	})
}

func (instance *HttpContext) AssertNamedHttpRequestResponseExactHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseExactHeaders(table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseContainsHeaders(table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponseContainsHeaders(alias string, table *godog.Table) error {
	return nil
}

func (instance *HttpContext) AssertResponseHeader(name, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseContentType(value string) error {
	return instance.AssertResponseHeader("content-type", value)
}

func (instance *HttpContext) AssertNamedHttpRequestResponseContentType(alias, value string) error {
	return instance.AssertNamedHttpRequestResponseHeader(alias, "content-type", value)
}

func (instance *HttpContext) AssertNamedHttpRequestResponseHeader(alias string, header string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponseIsValidAgainstSchema(file string) error {
	return instance.AssertNamedHttpRequestResponseIsValidAgainstSchema("", file)
}

func (instance *HttpContext) onNamedHttpRequestResponseWithJsonContentType(alias string, f func(*HttpRequest, *HttpResponse) error) error {
	return instance.onNamedResponse(alias, func(req *HttpRequest, res *HttpResponse) error {
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

func (instance *HttpContext) AssertNamedHttpRequestResponseIsValidAgainstSchema(alias, value string) error {
	return instance.onNamedHttpRequestResponseWithJsonContentType(alias, func(_ *HttpRequest, response *HttpResponse) error {
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

func (instance *HttpContext) onNamedHttpRequestResponseBodyPath(t, alias string, f func(*HttpRequest, *HttpResponse, any) error) error {
	return instance.onNamedHttpRequestResponseWithJsonContentType(alias, func(req *HttpRequest, res *HttpResponse) error {
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

func (instance *HttpContext) AssertResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSame(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSameOrAfter(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsSameOrBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsBefore(k, v string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathLengthIs(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathLengthIs(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsGreaterThan(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsGreaterThanOrEqualTo(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsLesserThan(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsLesserThanOrEqualTo(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsGreaterThanValue(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsGreaterThanOrEqualToValue(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsLesserThanValue(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsLesserThanOrEqualToValue(k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsGreaterThan(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualTo(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsLesserThan(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsLesserThanOrEqualTo(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsGreaterThanValue(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsGreaterThanOrEqualToValue(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsLesserThanValue(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsLesserThanOrEqualToValue(alias, k string, value float64) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsInStringArray(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertResponsePathIsInNumericArray(k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsInStringArray(alias, k string, value string) error {
	return nil
}

func (instance *HttpContext) AssertNamedHttpRequestResponsePathIsInNumericArray(alias, k string, value string) error {
	return nil
}

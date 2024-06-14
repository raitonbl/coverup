package http

import (
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"strings"
)

type URIScheme string

const (
	noneUriScheme  URIScheme = ""
	fileUriScheme  URIScheme = "file"
	httpUriScheme  URIScheme = "http"
	httpsUriScheme URIScheme = "https"
)

func execOnJsonContentType(instance *HttpContext, alias string, h func(*HttpRequest, *HttpResponse) error) error {
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
		return h(req, res)
	})
}

func execOnJsonPath(instance *HttpContext, alias, t string, h func(*HttpRequest, *HttpResponse, any) error) error {
	return execOnJsonContentType(instance, alias, func(req *HttpRequest, res *HttpResponse) error {
		expr := t
		if strings.HasPrefix(expr, ".") {
			expr = expr[1:]
		}
		if len(res.body) == 0 {
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
		return h(req, res, valueOf)
	})
}

func newJsonSchemaValidator(instance *HttpContext, opts HandlerOpts) any {
	f := func(alias string, value string) error {
		return execOnJsonContentType(instance, alias, func(req *HttpRequest, res *HttpResponse) error {
			if opts.scheme == noneUriScheme {
				return fmt.Errorf("URI scheme must be defined")
			}
			if opts.scheme != fileUriScheme && opts.scheme != httpUriScheme && opts.scheme != httpsUriScheme {
				return fmt.Errorf(`URI scheme "%s" isn't supported`, opts.scheme)
			}
			if instance.schemas == nil {
				instance.schemas = make(map[string]any)
			}
			valueOf, hasValue := instance.schemas[value]
			if !hasValue {
				binary, err := doGetFromURI(instance, opts.scheme, value)
				if err != nil {
					return err
				}
				l := gojsonschema.NewBytesLoader(binary)
				instance.schemas[value] = l
				valueOf = l
			}
			schemaLoader := valueOf.(gojsonschema.JSONLoader)
			documentLoader := gojsonschema.NewBytesLoader(res.body)
			r, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				return err
			}
			if opts.isAffirmationExpected {
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
			return fmt.Errorf("schema respects the schema %s://%s when it shouldn't", opts.scheme, value)
		})
	}
	if opts.isAliasedFunction {
		return f
	}
	return func(value string) error {
		return f("", value)
	}
}

func doGetFromURI(instance *HttpContext, scheme URIScheme, value string) ([]byte, error) {
	if scheme == fileUriScheme {
		return instance.ctx.GetFS().ReadFile(value)
	} else {
		url := string(scheme) + "://" + value
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		res, err := instance.ctx.GetResourcesHttpClient().Do(req)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(res.Body)
	}
}

package v3

import (
	"fmt"
	"strings"
)

func onJsonContentTypeDo(instance *HttpContext, alias string, h func(*HttpRequest, *HttpResponse) error) error {
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

func onJsonPathDo(instance *HttpContext, alias, t string, h func(*HttpRequest, *HttpResponse, any) error) error {
	return onJsonContentTypeDo(instance, alias, func(req *HttpRequest, res *HttpResponse) error {
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

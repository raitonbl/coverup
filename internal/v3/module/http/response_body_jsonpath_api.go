package http

import (
	"fmt"
	"regexp"
	"strings"
)

var isoDateRegexp, _ = regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
var isoTimeRegexp, _ = regexp.Compile(`^([01]\d|2[0-3]):([0-5]\d)(:[0-5]\d)?(\.\d{3})?$`)
var isoDateTimeRegexp, _ = regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])T([01]\d|2[0-3]):([0-5]\d):([0-5]\d)(\.\d{3})?(Z|([+-](0[0-9]|1[0-3]):[0-5]\d))?$`)

type HandlerOpts struct {
	isAffirmationExpected  bool
	isAliasedFunction      bool
	ignoreCaseIfApplicable bool
	attemptValueResolution bool
	scheme                 URIScheme
}

type HandlerFactory func(instance *Scenario, opts HandlerOpts) any

func newJsonPathIsTime(instance *Scenario, opts HandlerOpts) any {
	return newJsonPathIsTemporal(instance, opts, "Time", isoTimeRegexp)
}

func newJsonPathIsDate(instance *Scenario, opts HandlerOpts) any {
	return newJsonPathIsTemporal(instance, opts, "Date", isoDateRegexp)
}

func newJsonPathIsDateTime(instance *Scenario, opts HandlerOpts) any {
	return newJsonPathIsTemporal(instance, opts, "DateTime", isoDateTimeRegexp)
}

func newJsonPathIsTemporal(instance *Scenario, opts HandlerOpts, definedType string, regex *regexp.Regexp) any {
	f := func(expr, alias string) error {
		return instance.onNamedHttpRequestResponseBodyPath(expr, alias, func(_ *HttpRequest, res *HttpResponse, v any) error {
			value, isString := v.(string)
			if !isString {
				return fmt.Errorf("$%s should be a string but got %v", expr, v)
			}
			if regex.Match([]byte(value)) == opts.isAffirmationExpected {
				return nil
			}
			if opts.isAffirmationExpected {
				return fmt.Errorf("expected %s but got %v", definedType, value)
			}
			return fmt.Errorf("expected non %s but got %v", definedType, value)
		})
	}
	if opts.isAliasedFunction {
		return f
	}
	return func(expr string) error {
		return f(expr, "")
	}
}

func newJsonPathEqualsTo(instance *Scenario, opts HandlerOpts) any {
	f := newJsonPathEqualsToAnyHandler(instance, opts)
	if opts.isAliasedFunction {
		return func(expr, alias, compareTo string) error {
			return f.(func(string, string, any) error)(expr, alias, compareTo)
		}
	}
	return func(expr, compareTo string) error {
		return f.(func(string, any) error)(expr, compareTo)
	}
}

func newJsonPathEqualsToFloat64(instance *Scenario, opts HandlerOpts) any {
	f := newJsonPathEqualsToAnyHandler(instance, opts)
	if opts.isAliasedFunction {
		return func(expr, alias string, compareTo float64) error {
			return f.(func(string, string, any) error)(expr, alias, compareTo)
		}
	}
	return func(expr string, compareTo float64) error {
		return f.(func(string, any) error)(expr, compareTo)
	}
}

func newJsonPathEqualsToBooleanHandler(instance *Scenario, opts HandlerOpts) any {
	f := newJsonPathEqualsToAnyHandler(instance, opts)
	if opts.isAliasedFunction {
		return func(expr, alias string, compareTo string) error {
			valueOf := compareTo == "true"
			return f.(func(string, string, any) error)(expr, alias, valueOf)
		}
	}
	return func(expr string, compareTo string) error {
		valueOf := compareTo == "true"
		return f.(func(string, any) error)(expr, valueOf)
	}
}

func newJsonPathEqualsToAnyHandler(instance *Scenario, opts HandlerOpts) any {
	f := func(expr, alias string, compareTo any) error {
		return instance.onNamedHttpRequestResponseBodyPath(expr, alias, func(_ *HttpRequest, response *HttpResponse, value any) error {
			if (value == compareTo) == opts.isAffirmationExpected {
				return nil
			}
			condition := "must"
			if !opts.isAffirmationExpected {
				condition += "n't"
			}
			if alias == "" {
				return fmt.Errorf(`$%s=%v %s be equal to %v`, expr, value, condition, compareTo)
			}
			return fmt.Errorf(`%s.$%s=%v %s be equal to %v`, alias, expr, condition, value, compareTo)
		})
	}
	if opts.isAliasedFunction {
		return f
	}
	return func(expr string, compareTo any) error {
		return f(expr, "", compareTo)
	}
}

func newJsonPathContainsHandler(instance *Scenario, opts HandlerOpts) any {
	return newStringOperationJsonPathHandler(instance, "contain", opts, strings.Contains)
}

func newJsonPathStartsWithHandler(instance *Scenario, opts HandlerOpts) any {
	return newStringOperationJsonPathHandler(instance, "starts with", opts, strings.HasPrefix)
}

func newJsonPathPatternHandler(instance *Scenario, opts HandlerOpts) any {
	return newStringOperationJsonPathHandler(instance, "matches pattern", opts, func(fromResponse string, value string) bool {
		r, err := regexp.Compile(value)
		if err != nil {
			//TODO LOG ERROR
			return false
		}
		return r.Match([]byte(fromResponse))
	})
}

func newJsonPathEndsWithHandler(instance *Scenario, opts HandlerOpts) any {
	return newStringOperationJsonPathHandler(instance, "ends with", opts, strings.HasSuffix)
}

func newStringOperationJsonPathHandler(instance *Scenario, operation string, opts HandlerOpts, predicate func(string, string) bool) any {
	f := func(expr, alias string, c string) error {
		return doOnStringOperation(instance, operation, opts, alias, expr, c, predicate)
	}
	if opts.isAliasedFunction {
		return f
	}
	return func(expr string, compareTo string) error {
		return f(expr, "", compareTo)
	}
}

func doOnStringOperation(instance *Scenario, operation string, opts HandlerOpts, alias, expr, c string, predicate func(string, string) bool) error {
	return execOnJsonPath(instance, alias, expr, func(_ *HttpRequest, res *HttpResponse, value any) error {
		if value == nil {
			if alias == "" {
				return fmt.Errorf(`$%s mustn't be undefined`, expr)
			}
			return fmt.Errorf(`%s.$%s mustn't be undefined`, alias, expr)
		}
		valueOf, isText := value.(string)
		if !isText {
			if alias == "" {
				return fmt.Errorf(`$%s must be a string but got %v`, expr, value)
			}
			return fmt.Errorf(`%s.$%s must be a string but got %v`, alias, expr, value)
		}
		compareTo := c
		if opts.attemptValueResolution {
			v, err := instance.ctx.GetValue(c)
			if err != nil {
				return err
			}
			compareTo, isText = v.(string)
			if !isText {
				if alias == "" {
					return fmt.Errorf(`%s must be a string but got %v`, c, value)
				}
				return fmt.Errorf(`%s must be a string but got %v`, c, value)
			}
		}
		v1 := valueOf
		v2 := compareTo
		if opts.ignoreCaseIfApplicable {
			v1 = strings.ToUpper(v1)
			v2 = strings.ToUpper(v2)
		}
		r := predicate(v1, v2)
		if r == opts.isAffirmationExpected {
			return nil
		}
		condition := "must"
		if !opts.isAffirmationExpected {
			condition += "n't"
		}
		if alias == "" {
			return fmt.Errorf(`$%s=%v %s %s %v`, expr, value, condition, operation, compareTo)
		}
		return fmt.Errorf(`%s.$%s=%v %s %s %v`, alias, expr, condition, operation, value, compareTo)
	})
}

package http

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/checks"
	"github.com/thoas/go-funk"
	"regexp"
	"strconv"
	"strings"
)

const (
	anyStringRegex        = `"([^"]*)"`
	boolRegex             = `(true|false)`
	anyNumber             = `(-?\d+(\.\d+)?)`
	arrayRegex            = `\[(.*?)\]`
	resolvableStringRegex = `"` + valueRegex + `"`
)

var (
	boolRegexp, _      = regexp.Compile(boolRegex)
	anyNumberRegexp, _ = regexp.Compile(anyNumber)
	valueRegexp, _     = regexp.Compile(valueRegex)
	regexpCache        = make(map[string]*regexp.Regexp)
)

type PathOperationSettings struct {
	IgnoreCase  bool
	ValueRegexp string
}

type PathOperations struct {
	ConvertToNumberIfNecessary bool
	Line                       string
	ExpressionPattern          string
	PhraseFactory              func(string) []string
	AliasedPhraseFactory       func(string) []string
	ExtractFromResponse        func(response *Response, expr string) (any, error)
}

func (instance *PathOperations) New(ctx api.StepDefinitionContext) {
	arr := []func(api.StepDefinitionContext){
		instance.enableRegexSupport,
		instance.enableEqualsToSupport,
		instance.enableStartsWithSupport,
		instance.enableEndsWithSupport,
		instance.enableContainsSupport,
		instance.enableLesserThanSupport,
		instance.enableLesserOrEqualToSupport,
		instance.enableGreaterThanSupport,
		instance.enableGreaterOrEqualToSupport,
		instance.enableAnyOf,
	}
	for _, f := range arr {
		f(ctx)
	}
}

func (instance *PathOperations) enableAnyOf(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be any of", false, []string{valueRegex, arrayRegex}, instance.anyOfAssertionFactory)
}

func (instance *PathOperations) anyOfAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			valueOf, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			var arr any
			if options.Settings.ValueRegexp == valueRegex {
				arr, err = c.Resolve(v.(string))
			} else {
				arr, err = parseArray(c, v.(string))
			}
			if err != nil {
				return err
			}
			isTrue := funk.Contains(arr, valueOf)
			if isTrue == options.AssertTrue {
				return nil
			}
			if options.AssertTrue {
				return fmt.Errorf("%v isn't part of %v", valueOf, arr)
			}
			return fmt.Errorf("%v is part of %v", valueOf, arr)
		}
		return instance.createHandler(options, f)
	}
}

func parseArray(c api.ScenarioContext, value string) ([]any, error) {
	var arr = make([]any, 0)
	t := strings.Split(value, ",")
	for _, each := range t {
		if boolRegexp.MatchString(each) {
			arr = append(arr, each == "true")
		} else if strings.HasPrefix(each, `"`) && strings.HasSuffix(each, `"`) {
			arr = append(arr, each[1:len(each)-1])
		} else if anyNumberRegexp.MatchString(each) {
			v, err := strconv.ParseFloat(each, 64)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		} else if valueRegexp.MatchString(each) {
			v, err := c.Resolve(each)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		} else {
			return nil, fmt.Errorf("%v isn't support as part of [%v]", each, value)
		}
	}
	return arr, nil
}

func (instance *PathOperations) enableLesserThanSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be lesser than", false, instance.getDefaultArgs(), func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
		return instance.numericComparisonAssertionFactory(options, checks.IsLesserThan)
	})
}

func (instance *PathOperations) enableLesserOrEqualToSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be lesser or equal to", false, instance.getDefaultArgs(), func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
		return instance.numericComparisonAssertionFactory(options, checks.IsLesserOrEqualTo)
	})
}

func (instance *PathOperations) enableGreaterThanSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be greater than", false, instance.getDefaultArgs(), func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
		return instance.numericComparisonAssertionFactory(options, checks.IsGreaterThan)
	})
}

func (instance *PathOperations) enableGreaterOrEqualToSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be greater or equal to", false, instance.getDefaultArgs(), func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
		return instance.numericComparisonAssertionFactory(options, checks.IsGreaterOrEqualTo)
	})
}

func (instance *PathOperations) numericComparisonAssertionFactory(options FactoryOpts[PathOperationSettings], predicate func(float64, float64) bool) api.HandlerFactory {
	parseFloat64 := func(attr string, value any) (float64, error) {
		if !checks.IsString(value) && !instance.ConvertToNumberIfNecessary {
			return 0, fmt.Errorf("%s must be a float64", attr)
		}
		if o, err := strconv.ParseFloat(value.(string), 64); err != nil {
			return 0, fmt.Errorf("cannot convert %v, from %s, into float64 due to:\n%v", value, attr, err)
		} else {
			return o, nil
		}
	}
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			compareTo := v
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			if addr, isString := compareTo.(string); isString && options.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(addr)
				if prob != nil {
					return prob
				}
				compareTo = valueOf
				if compareTo, err = parseFloat64(v.(string), compareTo); err != nil {
					return err
				}
			}
			fromResponse, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			if fromResponse, err = parseFloat64(expr, fromResponse); err != nil {
				return err
			}
			r := predicate(fromResponse.(float64), compareTo.(float64))
			if r == options.AssertTrue {
				return nil
			}
			return fmt.Errorf("response.%s[%s]=%v", instance.Line, expr, fromResponse)
		}
		return instance.createHandler(options, f)
	}
}

func (instance *PathOperations) enableRegexSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "match pattern", true, instance.getDefaultArgs(), instance.regexAssertionFactory)
}

func (instance *PathOperations) enableEqualsToSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "be equal to", true, instance.getDefaultArgs(), instance.equalsToAssertionFactory)
}

func (instance *PathOperations) enableStartsWithSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "start with", true, instance.getDefaultArgs(), instance.startsWithAssertionFactory)
}

func (instance *PathOperations) enableEndsWithSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "end with", true, instance.getDefaultArgs(), instance.endsWithAssertionFactory)
}

func (instance *PathOperations) enableContainsSupport(ctx api.StepDefinitionContext) {
	instance.enabledSupportTo(ctx, "contain", true, instance.getDefaultArgs(), instance.containsAssertionFactory)
}

func (instance *PathOperations) startsWithAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.HasPrefix)
}

func (instance *PathOperations) endsWithAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.HasSuffix)
}

func (instance *PathOperations) containsAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, strings.Contains)
}

func (instance *PathOperations) regexAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return instance.stringOperationAssertionFactory(options, func(value string, pattern string) bool {
		r, hasValue := regexpCache[pattern]
		if !hasValue {
			c, err := regexp.Compile(pattern)
			if err != nil {
				// TODO: LOG
				return false
			}
			r = c
			regexpCache[pattern] = r
		}
		return r.MatchString(value)
	})
}

func (instance *PathOperations) stringOperationAssertionFactory(options FactoryOpts[PathOperationSettings], predicate func(string, string) bool) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			compareTo := v
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			if addr, isString := compareTo.(string); isString && options.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(addr)
				if prob != nil {
					return prob
				}
				compareTo = valueOf
			}
			fromResponse, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			if !checks.IsString(fromResponse) {
				return fmt.Errorf("response.%s[%s]=%v should be a string", instance.Line, expr, fromResponse)
			}
			if !checks.IsString(compareTo) {
				return fmt.Errorf("argument should be a string")
			}
			r := false
			if options.Settings.IgnoreCase {
				r = predicate(strings.ToUpper(fromResponse.(string)), strings.ToUpper(compareTo.(string)))
			} else {
				r = predicate(fromResponse.(string), compareTo.(string))
			}
			if r == options.AssertTrue {
				return nil
			}
			return fmt.Errorf("response.%s[%s]=%v", instance.Line, expr, fromResponse)
		}
		return instance.createHandler(options, f)
	}
}

func (instance *PathOperations) equalsToAssertionFactory(options FactoryOpts[PathOperationSettings]) api.HandlerFactory {
	return func(c api.ScenarioContext) any {
		f := func(alias string, expr string, v any) error {
			compareTo := v
			res, err := getHttpResponse(c, alias)
			if err != nil {
				return err
			}
			if addr, isString := compareTo.(string); isString && options.ResolveValueBeforeAssertion {
				valueOf, prob := c.Resolve(addr)
				if prob != nil {
					return prob
				}
				compareTo = valueOf
			}
			valueOf, err := instance.ExtractFromResponse(res, expr)
			if err != nil {
				return err
			}
			isEqualTo := false
			if options.Settings.IgnoreCase && checks.IsString(valueOf) && checks.IsString(compareTo) {
				isEqualTo = strings.ToUpper(valueOf.(string)) == strings.ToUpper(compareTo.(string))
			} else {
				isEqualTo = valueOf == compareTo
			}
			if isEqualTo == options.AssertTrue {
				return nil
			}
			return fmt.Errorf("response.%s[%s]=%v", instance.Line, expr, valueOf)
		}
		return instance.createHandler(options, f)
	}
}

func (instance *PathOperations) enabledSupportTo(ctx api.StepDefinitionContext, operation string, allowsIgnoreCase bool, args []string, f func(options FactoryOpts[PathOperationSettings]) api.HandlerFactory) {
	verbs := []string{"should", "shouldn't"}
	aliases := []string{"", httpRequestRegex}
	//	args := []string{anyStringRegex, resolvableStringRegex, valueRegex, anyNumber, boolRegex}
	for _, verb := range verbs {
		step := api.StepDefinition{
			Description: fmt.Sprintf("Asserts that a specific %s response %s %s the specified value", instance.Line, fmt.Sprintf("%s %s", verb, operation), ComponentType),
			Options:     make([]api.Option, 0),
		}
		for _, alias := range aliases {
			for _, arg := range args {
				numberOfOptions := 1
				supportsIgnoreCase := allowsIgnoreCase && (arg == anyStringRegex || arg == resolvableStringRegex || arg == valueRegex)
				if supportsIgnoreCase {
					numberOfOptions = 2
				}
				for i := 0; i < numberOfOptions; i++ {
					isIgnoreCase := i == 1
					var phrases []string
					format := fmt.Sprintf(`%s %s %s %s %s`, instance.Line, instance.ExpressionPattern, verb, operation, arg)
					if alias == aliases[0] {
						phrases = instance.PhraseFactory(format)
					} else {
						phrases = instance.AliasedPhraseFactory(format)
					}
					for _, p := range phrases {
						phrase := p
						if isIgnoreCase {
							phrase += ", ignoring case"
						}
						phrase += "$"
						options := FactoryOpts[PathOperationSettings]{
							Settings: &PathOperationSettings{
								ValueRegexp: arg,
								IgnoreCase:  isIgnoreCase,
							},
							AssertTrue:                  verb == verbs[0],
							AssertAlias:                 alias == aliases[1],
							ResolveValueBeforeAssertion: arg != anyNumber && arg != boolRegex,
						}
						step.Options = append(step.Options, api.Option{
							Regexp:         phrase,
							Description:    step.Description,
							HandlerFactory: f(options),
						})
					}
				}
			}
		}
		ctx.Then(step)
	}
}

func (instance *PathOperations) createHandler(options FactoryOpts[PathOperationSettings], f any) any {
	if !options.AssertAlias {
		alias := ""
		if options.Settings.ValueRegexp == anyNumber {
			return func(expr string, value float64) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		} else if options.Settings.ValueRegexp == boolRegex {
			return func(expr string, value bool) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		} else {
			return func(expr, value string) error {
				return f.(func(string, string, any) error)(alias, expr, value)
			}
		}
	}
	if options.Settings.ValueRegexp == anyNumber {
		return func(alias string, expr string, value float64) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	} else if options.Settings.ValueRegexp == boolRegex {
		return func(alias string, expr string, value bool) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	} else {
		return func(alias string, expr, value string) error {
			return f.(func(string, string, any) error)(alias, expr, value)
		}
	}
}

func (instance *PathOperations) getDefaultArgs() []string {
	return []string{anyStringRegex, resolvableStringRegex, valueRegex, anyNumber, boolRegex}
}

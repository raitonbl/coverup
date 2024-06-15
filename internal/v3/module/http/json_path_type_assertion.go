package http

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/checks"
)

type JsonPathTypeAssertionStepFactory struct {
}

func (instance *JsonPathTypeAssertionStepFactory) New(ctx api.StepDefinitionContext, c *Scenario) {
	params := map[string]func(FactoryOpts[any]) api.HandlerFactory{
		// Basic
		"Number":  instance.createFactory(c, "Number", checks.IsFloat64),
		"String":  instance.createFactory(c, "String", checks.IsString),
		"Boolean": instance.createFactory(c, "Boolean", checks.IsBoolean),
		// Temporal
		"Time":     instance.createFactory(c, "Time", checks.IsTime),
		"Date":     instance.createFactory(c, "Date", checks.IsDate),
		"DateTime": instance.createFactory(c, "DateTime", checks.ISDateTime),
	}
	opts := make([]api.Option, 0)
	for dataType, h := range params {
		representations := createResponseBodyJsonPathRegexp(fmt.Sprintf("is %s", dataType))
		for _, r := range representations {
			opts = append(opts, api.Option{
				HandlerFactory: h(FactoryOpts[any]{
					AssertTrue:                  true,
					AssertAlias:                 false,
					ResolveValueBeforeAssertion: true,
				}),
				Regexp:      r,
				Description: fmt.Sprintf("Asserts that, the current HttpResponse, response path is represents a %s", dataType),
			})
		}
		representations = createAliasedRequestResponseBodyJsonPathRegexp(fmt.Sprintf("is %s", dataType))
		for _, r := range representations {
			opts = append(opts, api.Option{
				HandlerFactory: h(FactoryOpts[any]{
					AssertTrue:                  true,
					AssertAlias:                 true,
					ResolveValueBeforeAssertion: true,
				}),
				Regexp:      r,
				Description: fmt.Sprintf("Asserts that, a given HttpRequest, response path is a %s", dataType),
			})
		}
	}
	ctx.Then(api.StepDefinition{
		Options:     opts,
		Type:        "Then",
		Description: fmt.Sprintf("Asserts that an HttpRequest response path is of a specific type"),
	})
}

func (instance *JsonPathTypeAssertionStepFactory) createFactory(c *Scenario, dataType string, predicate func(any) bool) func(FactoryOpts[any]) api.HandlerFactory {
	return func(opts FactoryOpts[any]) api.HandlerFactory {
		if !opts.AssertAlias {
			return func(_ api.ScenarioContext) any {
				return func(alias, expr string) error {
					return instance.execute(c, dataType, predicate, opts, "", expr)
				}
			}
		}
		return func(_ api.ScenarioContext) any {
			return func(expr string) error {
				return instance.execute(c, dataType, predicate, opts, "", expr)
			}
		}
	}
}

func (instance *JsonPathTypeAssertionStepFactory) execute(c *Scenario, dataType string, predicate func(any) bool, opts FactoryOpts[any], alias, expr string) error {
	valueOf, err := c.getResponseValueFromExpression(alias, expr)
	if err != nil {
		return err
	}
	assertTrue := predicate(valueOf)
	if assertTrue == opts.AssertTrue {
		return nil
	}
	if opts.AssertTrue {
		return fmt.Errorf("%v isn't %s", valueOf, dataType)
	}
	return fmt.Errorf("%v shouldn't be %s", valueOf, dataType)
}

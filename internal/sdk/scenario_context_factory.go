package sdk

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	"io/fs"
)

type ScenarioContextFactory struct {
	Properties         []string
	FileSystem         fs.ReadFileFS
	steps              []api.StepDefinition
	Entities           map[string]api.Entity
	OnScenarioCreation func(*DefaultScenarioContext)
}

func (instance *ScenarioContextFactory) Step(definition api.StepDefinition) {
	instance.doStep("Step", definition)
}

func (instance *ScenarioContextFactory) Given(definition api.StepDefinition) {
	instance.doStep("Given", definition)
}

func (instance *ScenarioContextFactory) When(definition api.StepDefinition) {
	instance.doStep("When", definition)
}

func (instance *ScenarioContextFactory) Then(definition api.StepDefinition) {
	instance.doStep("Then", definition)
}

func (instance *ScenarioContextFactory) Configure(c *godog.ScenarioContext) error {
	if instance.steps == nil {
		instance.steps = make([]api.StepDefinition, 0)
	}
	if instance.Properties == nil {
		instance.Properties = make([]string, 0)
	}
	sc := &DefaultScenarioContext{
		Filesystem: instance.FileSystem,
		Vars:       make(map[string]any),
		References: make(map[string]api.Component),
		Resolvers:  make(map[string]ValueResolver),
		Aliases:    make(map[string]map[string]api.Component),
	}
	if instance.OnScenarioCreation != nil {
		instance.OnScenarioCreation(sc)
	}
	// Assure fs from the current context is passed downstream
	sc.Filesystem = instance.FileSystem
	// Assure entities from the current context is passed downstream
	if instance.Entities != nil {
		for k, v := range instance.Entities {
			_ = sc.doAddGivenComponent(api.EntityComponentType, v, k, true)
		}
	}
	if sc.Resolvers == nil {
		sc.Resolvers = make(map[string]ValueResolver)
	}
	// Assure the original resolvers are passed down
	engine, err := NewPropertiesEngine(sc.Filesystem, instance.Properties...)
	if err != nil {
		return err
	}
	sc.Resolvers[api.PropertiesComponentType] = engine
	
	for _, definition := range instance.steps {
		for _, option := range definition.Options {
			if definition.Type == "Given" {
				c.Given(option.Regexp, option.HandlerFactory(sc))
			} else if definition.Type == "When" {
				c.When(option.Regexp, option.HandlerFactory(sc))
			} else if definition.Type == "Then" {
				c.Then(option.Regexp, option.HandlerFactory(sc))
			} else {
				c.Step(option.Regexp, option.HandlerFactory(sc))
			}
		}
	}
	return nil
}

func (instance *ScenarioContextFactory) doStep(stepType string, definition api.StepDefinition) {
	if instance.steps == nil {
		instance.steps = make([]api.StepDefinition, 0)
	}
	if definition.Description == "" {
		panic("step must have a description")
	}
	if definition.Options == nil || len(definition.Options) == 0 {
		return
	}
	for _, option := range definition.Options {
		if option.HandlerFactory == nil {
			panic("option must have a handler")
		}
		if option.Regexp == "" {
			panic("option must have a regexp")
		}
		if option.Description == "" {
			panic("option must have a description")
		}
	}
	instance.steps = append(instance.steps, api.StepDefinition{
		Type:        stepType,
		Options:     definition.Options,
		Description: definition.Description,
	})
}

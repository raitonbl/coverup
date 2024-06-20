package sdk

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/api/entities"
	"io/fs"
)

type ScenarioDefinitionContext struct {
	FileSystem         fs.ReadFileFS
	steps              []api.StepDefinition
	Entities           map[string]entities.Entity
	OnScenarioCreation func(*DefaultScenarioContext)
}

func (instance *ScenarioDefinitionContext) Step(definition api.StepDefinition) {
	instance.doStep("Step", definition)
}

func (instance *ScenarioDefinitionContext) Given(definition api.StepDefinition) {
	instance.doStep("Given", definition)
}

func (instance *ScenarioDefinitionContext) When(definition api.StepDefinition) {
	instance.doStep("When", definition)
}

func (instance *ScenarioDefinitionContext) Then(definition api.StepDefinition) {
	instance.doStep("Then", definition)
}

func (instance *ScenarioDefinitionContext) Configure(c *godog.ScenarioContext) {
	if instance.steps == nil {
		instance.steps = make([]api.StepDefinition, 0)
	}
	sc := &DefaultScenarioContext{
		Filesystem: instance.FileSystem,
		Vars:       make(map[string]any),
		References: make(map[string]api.Component),
		Entities:   make(map[string]entities.Entity),
		Aliases:    make(map[string]map[string]api.Component),
	}
	if instance.OnScenarioCreation != nil {
		instance.OnScenarioCreation(sc)
	}
	// Assure fs from the current context is passed downstream
	sc.Filesystem = instance.FileSystem
	// Assure entities from the current context is passed downstream
	if instance.Entities != nil {
		if sc.Entities == nil {
			sc.Entities = make(map[string]entities.Entity)
		}
		for id, each := range instance.Entities {
			sc.Entities[id] = each
		}
	}
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
}

func (instance *ScenarioDefinitionContext) doStep(stepType string, definition api.StepDefinition) {
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

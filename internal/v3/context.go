package v3

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
)

type ScenarioContext interface {
	context.Context
	GetValue(value string) (any, error)
	GerkhinContext() *godog.ScenarioContext
	GetComponent(componentType, alias string) (any, error)
	Register(componentType string, ptr context.Component, alias string) error
}

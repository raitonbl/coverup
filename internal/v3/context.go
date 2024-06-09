package v3

import (
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/internal/context"
	"github.com/raitonbl/coverup/pkg"
)

const valueExpression = `\{\{\s*([a-zA-Z0-9_]+\.)*[a-zA-Z0-9_]+\s*\}\}`

type ScenarioContext interface {
	context.Context
	GetValue(value string) (any, error)
	GerkhinContext() *godog.ScenarioContext
	GetComponent(componentType, alias string) (any, error)
	Register(componentType string, ptr pkg.Component, alias string) error
}

package pkg

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg/api"
	"io/fs"
)

const ValueExpression = `\{\{\s*([a-zA-Z0-9_]+\.)*[a-zA-Z0-9_]+\s*\}\}`

type ScenarioContext interface {
	context.Context
	GetFS() fs.ReadFileFS
	GetValue(value string) (any, error)
	GerkhinContext() *godog.ScenarioContext
	GetComponent(componentType, alias string) (any, error)
	Register(componentType string, ptr api.Component, alias string) error
}

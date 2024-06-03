package context

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg"
)

type Builder struct {
	context    Context
	references map[string]any
	aliases    map[string]map[string]any
}

func New(ctx Context) *Builder {
	return &Builder{
		context:    ctx,
		references: make(map[string]any),
		aliases:    make(map[string]map[string]any),
	}
}

func (instance *Builder) WithComponent(componentType string, ptr any, alias string) error {
	if alias != "" {
		_, hasValue := instance.aliases[componentType][alias]
		if hasValue {
			return fmt.Errorf("%s with alias %s cannot be defined more than once", componentType, alias)
		}
		instance.aliases[componentType][alias] = ptr
	}
	instance.references[componentType] = ptr
	return nil
}

func (instance *Builder) GetComponent(componentType string, alias string) any {
	if alias != "" {
		return instance.aliases[componentType][alias]
	}
	return instance.references[componentType]
}

// TODO: ALIGN WITH CONTEXT

func (instance *Builder) GetServerURL() string {
	return instance.context.GetServerURL()
}

func (instance *Builder) GetWorkDirectory() string {
	return instance.context.GetWorkDirectory()
}

func (instance *Builder) GetHttpClient() pkg.HttpClient {
	return instance.context.GetHttpClient()
}

func (instance *Builder) GetResourcesHttpClient() pkg.HttpClient {
	return instance.context.GetResourcesHttpClient()
}

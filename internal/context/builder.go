package context

import (
	"errors"
	"fmt"
	"github.com/raitonbl/coverup/pkg"
	"regexp"
)

var valueRegexp = regexp.MustCompile(`\{\{([^.]+)\.([^.]+)\.(.+?)\}\}`)

type Builder struct {
	context    Context
	references map[string]Component
	aliases    map[string]map[string]Component
}

func New(ctx Context) *Builder {
	return &Builder{
		context:    ctx,
		references: make(map[string]Component),
		aliases:    make(map[string]map[string]Component),
	}
}

func (instance *Builder) GetComponent(componentType string, alias string) Component {
	if alias != "" {
		if components, exists := instance.aliases[componentType]; exists {
			return components[alias]
		}
		return nil
	}
	return instance.references[componentType]
}

func (instance *Builder) WithComponent(componentType string, ptr Component, alias string) error {
	if alias != "" {
		if _, hasValue := instance.aliases[componentType]; !hasValue {
			instance.aliases[componentType] = make(map[string]Component)
		}
		if _, hasValue := instance.aliases[componentType][alias]; hasValue {
			return fmt.Errorf("%s with alias %s cannot be defined more than once", componentType, alias)
		}
		instance.aliases[componentType][alias] = ptr
	} else {
		instance.references[componentType] = ptr
	}
	return nil
}

func (instance *Builder) ResolveOrGetValue(p string) (any, error) {
	var prob error
	r := valueRegexp.ReplaceAllStringFunc(p, func(currentValue string) string {
		s := valueRegexp.FindStringSubmatch(currentValue)
		if len(s) == 0 {
			return currentValue
		}
		if len(s) != 4 {
			prob = errors.New("cannot resolve " + p)
			return ""
		}
		// Extract the captured groups
		componentType := s[1]
		componentId := s[2]
		path := s[3]

		components, hasValue := instance.aliases[componentType]
		if !hasValue {
			prob = errors.New("Component " + componentType + "." + componentId + " has not been defined yet")
			return ""
		}
		component, hasValue := components[componentId]
		if !hasValue {
			prob = errors.New("Component " + componentType + "." + componentId + " has not been defined yet")
			return ""
		}
		valueOf, err := component.GetPathValue(path)
		if err != nil {
			prob = err
			return ""
		}
		return fmt.Sprintf("%v", valueOf)
	})
	if prob != nil {
		return nil, prob
	}
	return r, nil
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

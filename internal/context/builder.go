package context

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg"
	"regexp"
	"strings"
)

var valueRegexp = regexp.MustCompile(`{{\s*([\w.]+)\s*}}`)

type Builder struct {
	context    Context
	references map[string]pkg.Component
	aliases    map[string]map[string]pkg.Component
}

func New(ctx Context) *Builder {
	return &Builder{
		context:    ctx,
		references: make(map[string]pkg.Component),
		aliases:    make(map[string]map[string]pkg.Component),
	}
}

func (instance *Builder) GetComponent(componentType string, alias string) pkg.Component {
	if alias != "" {
		if components, exists := instance.aliases[componentType]; exists {
			return components[alias]
		}
		return nil
	}
	return instance.references[componentType]
}

func (instance *Builder) WithComponent(componentType string, ptr pkg.Component, alias string) error {
	if alias != "" {
		if _, hasValue := instance.aliases[componentType]; !hasValue {
			instance.aliases[componentType] = make(map[string]pkg.Component)
		}
		if _, hasValue := instance.aliases[componentType][alias]; hasValue {
			return fmt.Errorf("%s with alias %s cannot be defined more than once", componentType, alias)
		}
		instance.aliases[componentType][alias] = ptr
	}
	instance.references[componentType] = ptr
	return nil
}
func (instance *Builder) GetValue(src string) (any, error) {
	matches := valueRegexp.FindAllStringSubmatch(src, -1)
	if len(matches) == 0 {
		return src, nil
	}
	cache := make(map[string]string)
	parsedValue := src
	for _, match := range matches {
		if len(match) <= 1 {
			continue
		}
		key := match[1]
		expr := match[0]
		v, err := instance.getValueFromExpression(key, expr, cache)
		if err != nil {
			return nil, err
		}
		parsedValue = strings.ReplaceAll(parsedValue, expr, v)
	}

	return parsedValue, nil
}

func (instance *Builder) getValueFromExpression(key, expr string, cache map[string]string) (string, error) {
	if fromCache, containsKey := cache[key]; containsKey {
		return fromCache, nil
	}
	arr := strings.Split(key, ".")
	if len(arr) < 3 {
		return "", fmt.Errorf(`expression %s is unresolvable`, expr)
	}
	componentType, componentId, path := arr[0], arr[1], strings.Join(arr[2:], ".")
	component, err := instance.getComponentOrElseThrow(componentType, componentId)
	if err != nil {
		return "", err
	}
	valueOf, err := component.GetPathValue(path)
	if err != nil {
		return "", err
	}
	v := fmt.Sprintf("%v", valueOf)
	cache[key] = v

	return v, nil
}

func (instance *Builder) getComponentOrElseThrow(componentType, componentId string) (pkg.Component, error) {
	components, hasValue := instance.aliases[componentType]
	if !hasValue {
		return nil, fmt.Errorf("component %s.%s has not been defined yet", componentType, componentId)
	}

	component, hasValue := components[componentId]
	if !hasValue {
		return nil, fmt.Errorf("component %s.%s has not been defined yet", componentType, componentId)
	}

	return component, nil
}

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

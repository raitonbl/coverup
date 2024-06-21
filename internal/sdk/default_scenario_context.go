package sdk

import (
	"fmt"
	"github.com/raitonbl/coverup/pkg/api"
	"io/fs"
	"regexp"
	"strings"
)

var valueRegexp = regexp.MustCompile(api.ValueExpression)

type DefaultScenarioContext struct {
	Filesystem fs.ReadFileFS
	Vars       map[string]any
	References map[string]api.Component
	Aliases    map[string]map[string]api.Component
	Resolvers  map[string]ValueResolver
}

func (d *DefaultScenarioContext) GetFS() fs.ReadFileFS {
	return d.Filesystem
}

func (d *DefaultScenarioContext) Resolve(src string) (any, error) {
	// [objectType].[objectId].<property>
	// DynamodbItem.<id>.<property>
	// HttpRequest.<name>.<property>
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
		v, err := d.getValueFromExpression(key, expr, cache)
		if err != nil {
			return nil, err
		}
		parsedValue = strings.ReplaceAll(parsedValue, expr, v)
	}
	return parsedValue, nil
}

func (d *DefaultScenarioContext) GetGivenComponent(componentType, alias string) (any, error) {
	if alias != "" {
		if components, exists := d.Aliases[componentType]; exists {
			return components[alias], nil
		}
		return nil, nil
	}
	return d.References[componentType], nil
}

func (d *DefaultScenarioContext) getValueFromExpression(key, expr string, cache map[string]string) (string, error) {
	if fromCache, containsKey := cache[key]; containsKey {
		return fromCache, nil
	}
	indexOf := strings.Index(expr, ".")
	componentType := expr[:indexOf]
	if d.Resolvers == nil {
		d.Resolvers = make(map[string]ValueResolver)
	}
	var err error
	var valueOf any
	r, hasValue := d.Resolvers[componentType]
	if hasValue {
		valueOf, err = r.ValueFrom(expr[indexOf+1:])
	} else {
		valueOf, err = d.ValueFrom(expr)
	}
	if err != nil {
		return "", err
	}
	v := fmt.Sprintf("%v", valueOf)
	cache[key] = v
	return v, nil
}

func (d *DefaultScenarioContext) ValueFrom(expr string) (any, error) {
	arr := strings.Split(expr, ".")
	if len(arr) < 2 {
		return "", fmt.Errorf(`expression %s is unresolvable`, expr)
	}
	componentType, componentId, path := arr[0], arr[1], strings.Join(arr[2:], ".")
	component, err := d.getComponentOrElseThrow(componentType, componentId)
	if err != nil {
		return "", err
	}
	return component.ValueFrom(path)
}

func (d *DefaultScenarioContext) getComponentOrElseThrow(componentType, componentId string) (api.Component, error) {
	components, hasValue := d.Aliases[componentType]
	if !hasValue {
		return nil, fmt.Errorf("component %s.%s has not been defined yet", componentType, componentId)
	}

	component, hasValue := components[componentId]
	if !hasValue {
		return nil, fmt.Errorf("component %s.%s has not been defined yet", componentType, componentId)
	}

	return component, nil
}

func (d *DefaultScenarioContext) AddGivenComponent(componentType string, ptr api.Component, alias string) error {
	if componentType == api.EntityComponentType {
		return fmt.Errorf("cannot add a component with type %s", componentType)
	}
	return d.doAddGivenComponent(componentType, ptr, alias, false)
}

func (d *DefaultScenarioContext) doAddGivenComponent(componentType string, ptr api.Component, alias string, allowOverride bool) error {
	if alias != "" {
		if _, hasValue := d.Aliases[componentType]; !hasValue {
			d.Aliases[componentType] = make(map[string]api.Component)
		}
		if _, hasValue := d.Aliases[componentType][alias]; hasValue && !allowOverride {
			return fmt.Errorf("%s with alias %s cannot be defined more than once", componentType, alias)
		}
		d.Aliases[componentType][alias] = ptr
	}
	d.References[componentType] = ptr
	return nil
}

func (d *DefaultScenarioContext) GetValue(namespace, key string) (any, error) {
	if d.Vars == nil {
		d.Vars = make(map[string]any)
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be undefined")
	}
	if key == "" {
		return nil, fmt.Errorf("key cannot be undefined")
	}
	return d.Vars[namespace+"."+key], nil
}

func (d *DefaultScenarioContext) SetValue(namespace, key string, value any) error {
	if d.Vars == nil {
		d.Vars = make(map[string]any)
	}
	if namespace == "" {
		return fmt.Errorf("namespace cannot be undefined")
	}
	if key == "" {
		return fmt.Errorf("key cannot be undefined")
	}
	d.Vars[namespace+"."+key] = value
	return nil
}

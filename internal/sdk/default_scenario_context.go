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
}

func (d *DefaultScenarioContext) GetFS() fs.ReadFileFS {
	return d.Filesystem
}

func (d *DefaultScenarioContext) Resolve(src string) (any, error) {
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
	arr := strings.Split(key, ".")
	if len(arr) < 3 {
		return "", fmt.Errorf(`expression %s is unresolvable`, expr)
	}
	componentType, componentId, path := arr[0], arr[1], strings.Join(arr[2:], ".")
	component, err := d.getComponentOrElseThrow(componentType, componentId)
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
	if alias != "" {
		if _, hasValue := d.Aliases[componentType]; !hasValue {
			d.Aliases[componentType] = make(map[string]api.Component)
		}
		if _, hasValue := d.Aliases[componentType][alias]; hasValue {
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

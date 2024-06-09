package v3

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/raitonbl/coverup/pkg"
	"regexp"
	"strings"
)

var valueRegexp = regexp.MustCompile(valueExpression)

type DefaultScenarioContext struct {
	WorkDirectory string
	HttpClient    pkg.HttpClient
	GoDogContext  *godog.ScenarioContext
	References    map[string]pkg.Component
	Aliases       map[string]map[string]pkg.Component
}

func (d *DefaultScenarioContext) GetServerURL() string {
	panic("implement me")
}

func (d *DefaultScenarioContext) GetWorkDirectory() string {
	return d.WorkDirectory
}

func (d *DefaultScenarioContext) GetHttpClient() pkg.HttpClient {
	return d.HttpClient
}

func (d *DefaultScenarioContext) GetResourcesHttpClient() pkg.HttpClient {
	return d.HttpClient
}

func (d *DefaultScenarioContext) GetValue(src string) (any, error) {
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

func (d *DefaultScenarioContext) getComponentOrElseThrow(componentType, componentId string) (pkg.Component, error) {
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

func (d *DefaultScenarioContext) GerkhinContext() *godog.ScenarioContext {
	return d.GoDogContext
}

func (d *DefaultScenarioContext) GetComponent(componentType, alias string) (any, error) {
	if alias != "" {
		if components, exists := d.Aliases[componentType]; exists {
			return components[alias], nil
		}
		return nil, nil
	}
	return d.References[componentType], nil
}

func (d *DefaultScenarioContext) Register(componentType string, ptr pkg.Component, alias string) error {
	if alias != "" {
		if _, hasValue := d.Aliases[componentType]; !hasValue {
			d.Aliases[componentType] = make(map[string]pkg.Component)
		}
		if _, hasValue := d.Aliases[componentType][alias]; hasValue {
			return fmt.Errorf("%s with alias %s cannot be defined more than once", componentType, alias)
		}
		d.Aliases[componentType][alias] = ptr
	}
	d.References[componentType] = ptr
	return nil
}

package internal

import (
	"fmt"
	"github.com/cucumber/godog"
	v3Pkg "github.com/raitonbl/coverup/internal/pkg"
	"github.com/raitonbl/coverup/pkg/api"
	"github.com/raitonbl/coverup/pkg/http"
	"io/fs"
	"regexp"
	"strings"
)

var valueRegexp = regexp.MustCompile(v3Pkg.ValueExpression)

type DefaultScenarioContext struct {
	Filesystem   fs.ReadFileFS
	HttpClient   http.Client
	GoDogContext *godog.ScenarioContext
	References   map[string]api.Component
	Aliases      map[string]map[string]api.Component
}

func (d *DefaultScenarioContext) GetServerURL() string {
	panic("implement me")
}

func (d *DefaultScenarioContext) GetHttpClient() http.Client {
	return d.HttpClient
}

func (d *DefaultScenarioContext) GetResourcesHttpClient() http.Client {
	return d.HttpClient
}

func (d *DefaultScenarioContext) GetFS() fs.ReadFileFS {
	return d.Filesystem
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

func (d *DefaultScenarioContext) Register(componentType string, ptr api.Component, alias string) error {
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

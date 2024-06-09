package v3

import (
	"embed"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/pkg"
	"os"
	"strings"
	"testing"
)

//go:embed testdata/*
var homeDirectory embed.FS

func TestApply(t *testing.T) {
	doApply(t, "features/design-api/default.feature", nil)
}

func doApply(t *testing.T, filename string, f func(int) error) {
	binary, err := homeDirectory.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workDir += "/testdata/features/design-api"
	status := godog.TestSuite{
		TestSuiteInitializer: nil,
		Options: &godog.Options{
			TestingT:      t,
			Format:        "pretty",
			StopOnFailure: true,
			Paths:         []string{},
			FeatureContents: []godog.Feature{
				{
					Contents: binary,
					Name:     t.Name(),
				},
			},
			Output: colors.Colored(os.Stdout),
			Strict: true,
		},
		ScenarioInitializer: func(gherkinContext *godog.ScenarioContext) {
			ctxt := &V3Context{
				workDirectory:  workDir,
				gherkinContext: gherkinContext,
				references:     make(map[string]pkg.Component),
				aliases:        make(map[string]map[string]pkg.Component),
			}
			On(ctxt)
		},
	}.Run()
	if f != nil {
		err = f(status)
	}
	if err != nil {
		t.Fatal(err)
	}

}

type V3Context struct {
	workDirectory  string
	gherkinContext *godog.ScenarioContext
	references     map[string]pkg.Component
	aliases        map[string]map[string]pkg.Component
}

func (instance *V3Context) GetServerURL() string {
	//TODO implement me
	panic("implement me")
}

func (instance *V3Context) GetWorkDirectory() string {
	if instance.workDirectory == "" {
		return "./"
	}
	return instance.workDirectory
}

func (instance *V3Context) GetHttpClient() pkg.HttpClient {
	//TODO implement me
	panic("implement me")
}

func (instance *V3Context) GetResourcesHttpClient() pkg.HttpClient {
	//TODO implement me
	panic("implement me")
}

func (instance *V3Context) GerkhinContext() *godog.ScenarioContext {
	return instance.gherkinContext
}

func (instance *V3Context) Register(componentType string, ptr pkg.Component, alias string) error {
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

func (instance *V3Context) GetValue(value string) (any, error) {
	if strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}") {
		return "picture.base64", nil
	}
	return value, nil
}

func (instance *V3Context) GetComponent(componentType, alias string) (any, error) {
	if alias != "" {
		if components, exists := instance.aliases[componentType]; exists {
			return components[alias], nil
		}
		return nil, nil
	}
	return instance.references[componentType], nil
}

package config

import (
	"bytes"
	"github.com/raitonbl/coverup/internal/sdk"
	"github.com/raitonbl/coverup/pkg/api"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"text/template"
)

func Load(ctx api.ScenarioContext) (*Manifest, error) {
	manifest := &Manifest{}
	binary, err := ctx.GetFS().ReadFile("aws.conf")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if binary != nil {
		tmpl, prob := template.New("config").Parse(string(binary))
		if prob != nil {
			return nil, prob
		}
		var bf bytes.Buffer
		if prob = tmpl.Execute(&bf, binary); err != nil {
			return nil, prob
		}
		templateData, prob := createTemplateData(ctx)
		if prob != nil {
			return nil, prob
		}
		if prob = yaml.Unmarshal(bf.Bytes(), templateData); prob != nil {
			return nil, prob
		}
	}
	if manifest.Config == nil {
		manifest.Config = make(map[string]Config)
	}
	_, hasDefault := manifest.Config["default"]
	defaultConfig := Config{
		Credentials: &Credentials{
			Type: "default",
		},
	}
	if !hasDefault {
		manifest.Config["default"] = defaultConfig
	}
	if manifest.Services != nil {
		for _, v := range manifest.Services {
			if v.Config == nil {
				v.Config = &defaultConfig
			}
		}
	}
	return manifest, nil
}

func createTemplateData(ctx api.ScenarioContext) (any, error) {
	component, err := ctx.GetGivenComponent(api.PropertiesComponentType, "")
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)
	data["Environment"] = getEnvironmentVariables()
	data["Properties"] = component.(*sdk.PropertiesEngine).ToMap()
	return data, nil
}

func getEnvironmentVariables() map[string]string {
	variables := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			variables[pair[0]] = pair[1]
		}
	}
	return variables
}

package v3

import (
	"embed"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/pkg"
	"os"
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
		},
		ScenarioInitializer: func(gherkinContext *godog.ScenarioContext) {
			ctxt := &V3Context{
				gherkinContext: gherkinContext,
			}
			Apply(ctxt)
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
	gherkinContext *godog.ScenarioContext
}

func (v *V3Context) GetServerURL() string {
	//TODO implement me
	panic("implement me")
}

func (v *V3Context) GetWorkDirectory() string {
	//TODO implement me
	panic("implement me")
}

func (v *V3Context) GetHttpClient() pkg.HttpClient {
	//TODO implement me
	panic("implement me")
}

func (v *V3Context) GetResourcesHttpClient() pkg.HttpClient {
	//TODO implement me
	panic("implement me")
}

func (v *V3Context) GerkhinContext() *godog.ScenarioContext {
	return v.gherkinContext
}

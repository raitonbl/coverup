package v3

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/pkg"
	"net/http"
	"os"
	"testing"
)

func Exec(t *testing.T, definition []byte, c map[string]func(*http.Request) (*http.Response, error), _ map[string][]byte) {
	workDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	httpClient := &FnHttpClient{
		c,
	}
	suite := godog.TestSuite{
		TestSuiteInitializer: nil,
		Options: &godog.Options{
			TestingT:      t,
			StopOnFailure: true,
			Format:        "pretty",
			Paths:         []string{},
			FeatureContents: []godog.Feature{{
				Contents: definition,
				Name:     t.Name(),
			},
			},
			Output: colors.Colored(os.Stdout),
		},
		ScenarioInitializer: func(goDogCtx *godog.ScenarioContext) {
			ctx := &DefaultScenarioContext{
				GoDogContext:  goDogCtx,
				HttpClient:    httpClient,
				WorkDirectory: workDirectory,
				References:    make(map[string]pkg.Component),
				Aliases:       make(map[string]map[string]pkg.Component),
			}
			On(ctx)
		},
	}
	suite.Run()
}

type FnHttpClient struct {
	m map[string]func(*http.Request) (*http.Response, error)
}

func (f *FnHttpClient) Do(req *http.Request) (*http.Response, error) {
	k := req.Method + " " + req.URL.String()
	return f.m[k](req)
}

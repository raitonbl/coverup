package v3

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/pkg"
	"io/fs"
	"net/http"
	"os"
	"testing"
)

func Exec(t *testing.T, definition []byte, c map[string]func(*http.Request) (*http.Response, error), fm map[string]func() ([]byte, error)) {
	filesystem := &FnFS{
		fm,
	}
	if fm != nil {
		filesystem.m = fm
	}
	httpClient := &FnHttpClient{
		c,
	}
	suite := godog.TestSuite{
		TestSuiteInitializer: nil,
		Options: &godog.Options{
			TestingT:      t,
			Strict:        true,
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
				Filesystem:   filesystem,
				GoDogContext: goDogCtx,
				HttpClient:   httpClient,
				References:   make(map[string]pkg.Component),
				Aliases:      make(map[string]map[string]pkg.Component),
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

type FnFS struct {
	m map[string]func() ([]byte, error)
}

func (f FnFS) Open(name string) (fs.File, error) {
	panic("implement me")
}

func (f FnFS) ReadFile(name string) ([]byte, error) {
	return f.m[name]()
}

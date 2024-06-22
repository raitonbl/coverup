package http

import (
	"bytes"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/internal/sdk"
	"github.com/raitonbl/coverup/pkg/api"
	pkgHttp "github.com/raitonbl/coverup/pkg/http"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"net/http"
	"os"
	"testing"
)

type GivenRequestOpts struct {
	properties []string
	filesystem fs.ReadFileFS
	httpClient pkgHttp.Client
	entities   map[string]api.Entity
}

func TestGivenRequestOnBehalfOfEntity(t *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	r := readProductFromFile(id)
	filesystem := &FnFS{
		map[string]func() ([]byte, error){
			"requests/product.json": func() ([]byte, error) {
				return r, nil
			},
			"schemas/product.json": func() ([]byte, error) {
				return getProductJSONSchema(), nil
			},
		},
	}
	httpClient := &FnHttpClient{
		m: map[string]func(*http.Request) (*http.Response, error){
			fmt.Sprintf("GET https://localhost:8443/items/%s", id): func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Status:     http.StatusText(200),
					Header: map[string][]string{
						"x-ratelimit-remaining": {"50"},
						"x-ratelimit-limit":     {"100"},
						"x-ratelimit-reset":     {"1625690400"},
						"content-type":          {"application/json"},
						"x-amzn-trace-id":       {"Root=1-5f84c3a3-91f49ffb0a2e26a3a3e58d0c; Parent=36b815b057b745d6; Sampled=1"},
					},
					Body: io.NopCloser(bytes.NewBuffer(r)),
				}, nil
			},
			"GET http://localhost:8080/schemas/product.json":  getProductSchemaFromURL,
			"GET https://localhost:8443/schemas/product.json": getProductSchemaFromURL,
		},
	}
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	ExecGivenRequest(t, []byte(`
		Feature: 
			Scenario:
			Given a HttpRequest
			And the http request headers are:
				| content-type | application/json |
			And the http request method is GET
			And http request URL is https://localhost:8443
			And http request path is /items/`+id+` 
			When {{`+api.EntityComponentType+`.default}} submits the HttpRequest
			Then http response status code should be 200`),
		GivenRequestOpts{
			filesystem: filesystem,
			httpClient: httpClient,
			entities: map[string]api.Entity{
				"default": api.BearerToken{
					BasicEntity: api.BasicEntity{
						Name:        "Bearer Token",
						Description: "Just a bearer token",
					},
					Value: bearerToken,
				},
			},
		})
	require.Equal(t, httpClient.data[0].Header.Get("Authorization"), fmt.Sprintf(`Bearer %s`, bearerToken))
}

func TestGivenRequestWhenAliased(t *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	r := readProductFromFile(id)
	filesystem := &FnFS{
		map[string]func() ([]byte, error){
			"requests/product.json": func() ([]byte, error) {
				return r, nil
			},
			"schemas/product.json": func() ([]byte, error) {
				return getProductJSONSchema(), nil
			},
		},
	}
	httpClient := &FnHttpClient{
		m: map[string]func(*http.Request) (*http.Response, error){
			fmt.Sprintf("GET https://localhost:8443/items/%s", id): func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Status:     http.StatusText(200),
					Header: map[string][]string{
						"x-ratelimit-remaining": {"50"},
						"x-ratelimit-limit":     {"100"},
						"x-ratelimit-reset":     {"1625690400"},
						"content-type":          {"application/json"},
						"x-amzn-trace-id":       {"Root=1-5f84c3a3-91f49ffb0a2e26a3a3e58d0c; Parent=36b815b057b745d6; Sampled=1"},
					},
					Body: io.NopCloser(bytes.NewBuffer(r)),
				}, nil
			},
			"GET http://localhost:8080/schemas/product.json":  getProductSchemaFromURL,
			"GET https://localhost:8443/schemas/product.json": getProductSchemaFromURL,
		},
	}
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	ExecGivenRequest(t, []byte(`
		Feature: 
			Scenario:
			Given a HttpRequest named GetItems
			And the http request headers are:
				| content-type | application/json |
			And the http request method is GET
			And http request URL is https://localhost:8443
			And http request path is /items/`+id+` 
			When {{`+api.EntityComponentType+`.default}} submits {{`+ComponentType+`.GetItems}}
			Then http response status code should be 200`),
		GivenRequestOpts{
			filesystem: filesystem,
			httpClient: httpClient,
			entities: map[string]api.Entity{
				"default": api.BearerToken{
					BasicEntity: api.BasicEntity{
						Name:        "Bearer Token",
						Description: "Just a bearer token",
					},
					Value: bearerToken,
				},
			},
		})
	require.Equal(t, httpClient.data[0].Header.Get("Authorization"), fmt.Sprintf(`Bearer %s`, bearerToken))
}

func TestGivenRequestWhenPropertiesURL(t *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	r := readProductFromFile(id)
	filesystem := &FnFS{
		map[string]func() ([]byte, error){
			"requests/product.json": func() ([]byte, error) {
				return r, nil
			},
			"schemas/product.json": func() ([]byte, error) {
				return getProductJSONSchema(), nil
			},
			"dev.properties": func() ([]byte, error) {
				return []byte(`server.url=https://localhost:8443`), nil
			},
		},
	}
	httpClient := &FnHttpClient{
		m: map[string]func(*http.Request) (*http.Response, error){
			fmt.Sprintf("GET https://localhost:8443/items/%s", id): func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Status:     http.StatusText(200),
					Header: map[string][]string{
						"x-ratelimit-remaining": {"50"},
						"x-ratelimit-limit":     {"100"},
						"x-ratelimit-reset":     {"1625690400"},
						"content-type":          {"application/json"},
						"x-amzn-trace-id":       {"Root=1-5f84c3a3-91f49ffb0a2e26a3a3e58d0c; Parent=36b815b057b745d6; Sampled=1"},
					},
					Body: io.NopCloser(bytes.NewBuffer(r)),
				}, nil
			},
			"GET http://localhost:8080/schemas/product.json":  getProductSchemaFromURL,
			"GET https://localhost:8443/schemas/product.json": getProductSchemaFromURL,
		},
	}
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	ExecGivenRequest(t, []byte(`
		Feature: 
			Scenario:
			Given a HttpRequest named GetItems
			And the http request headers are:
				| content-type | application/json |
			And the http request method is GET
			And http request URL is {{`+api.PropertiesComponentType+`.server.url}}
			And http request path is /items/`+id+` 
			When {{`+api.EntityComponentType+`.default}} submits {{`+ComponentType+`.GetItems}}
			Then http response status code should be 200`),
		GivenRequestOpts{
			filesystem: filesystem,
			httpClient: httpClient,
			properties: []string{"dev.properties"},
			entities: map[string]api.Entity{
				"default": api.BearerToken{
					BasicEntity: api.BasicEntity{
						Name:        "Bearer Token",
						Description: "Just a bearer token",
					},
					Value: bearerToken,
				},
			},
		})
}

func getProductSchemaFromURL(_ *http.Request) (*http.Response, error) {
	binary := getProductJSONSchema()
	return &http.Response{
		StatusCode: 200,
		Status:     http.StatusText(200),
		Header: map[string][]string{
			"content-type": {"application/json"},
		},
		Body: io.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}

func ExecGivenRequest(t *testing.T, definition []byte, givenRequestOpts GivenRequestOpts) {
	ctx := &sdk.ScenarioContextFactory{
		Entities:   givenRequestOpts.entities,
		FileSystem: givenRequestOpts.filesystem,
		Properties: givenRequestOpts.properties,
		OnScenarioCreation: func(context *sdk.DefaultScenarioContext) {
			_ = context.SetValue(ComponentType, "httpClient", givenRequestOpts.httpClient)
		},
	}
	OnV3(ctx)
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
			err := ctx.Configure(goDogCtx)
			if err != nil {
				t.Fatal(err)
			}
		},
	}
	suite.Run()
}

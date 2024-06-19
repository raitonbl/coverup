package http

import (
	"bytes"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/raitonbl/coverup/internal/v3/api"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestV3Api_on_body(m *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	doTestV3Api(m, id, []string{
		`http response body $.id should be equal to "` + id + `"`,
		`http response body $.id shouldn't be equal to "RaitonBL"`,
		`http response body $.summary should contain "SSD"`,
		`http response body $.summary should contain "ssd", ignoring case`,
		`http response body $.name should start with "Seagate"`,
		`http response body $.name should start with "SeaGate", ignoring case`,
		`http response body $.name should end with "MB/s"`,
		`http response body $.name should end with "mb/s", ignoring case`,
		`http response body $.warranty.amount should be equal to 2`,
		`http response body $.warranty.amount shouldn't be equal to 3`,
		`http response body $.in_promotion should be equal to false`,
		`http response body $.in_promotion shouldn't be equal to true`,
		`http response body $.name should match pattern "^Seagate"`,
		`http response body $.name shouldn't match pattern "^X"`,
		`http response body should be file://requests/product.json`,
		`http response body respects json schema file://schemas/product.json`,
		`http response body respects json schema http://localhost:8080/schemas/product.json`,
		`http response body respects json schema https://localhost:8443/schemas/product.json`,
		`http response body should be:
		"""
			` + string(readProductFromFile(id)) + `
		"""`,
	})
}

func TestV3Api_on_headers(m *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	doTestV3Api(m, id, []string{
		"http response status code should be 200",
		"http response status code shouldn't be 201",
		`http response headers should contain:
			| content-type | application/json |`,
		`http response headers should be:
			| content-type 	  | application/json 															|
			| x-amzn-trace-id | Root=1-5f84c3a3-91f49ffb0a2e26a3a3e58d0c; Parent=36b815b057b745d6; Sampled=1 |
			| x-ratelimit-remaining| 50|
			| x-ratelimit-limit | 100|
			| x-ratelimit-reset |1625690400|`,
		`http response headers shouldn't contain:
			| content-type | application/xml |`,
		`http response headers shouldn't be:
			| content-type 	  | application/problem+json 													|
			| x-amzn-trace-id | Root=1-5f84c3a3-91f49ffb0a2e26a3a3e58d0c; Parent=36b815b057b745d6; Sampled=1 |`,
		`http response body should be:
		"""
			` + string(readProductFromFile(id)) + `
		"""`,
		`http response header content-type should be equal to "application/json"`,
		`http response header content-type shouldn't be equal to "application/problem+json"`,
		`http response header content-type should be equal to "application/JSON", ignoring case`,
		`http response header content-type shouldn't be equal to "application/problem+json", ignoring case`,
		`http response header content-type should end with "json"`,
		`http response header content-type should start with "application"`,
		`http response header content-type should contain "/"`,
		`http response header content-type should match pattern "application/*"`,
		`http response header content-type shouldn't match pattern "*/xml"`,
		`http response header x-ratelimit-remaining should be greater than 49`,
		`http response header x-ratelimit-remaining should be greater or equal to 50`,
		`http response header x-ratelimit-limit should be lesser than 101`,
		`http response header x-ratelimit-remaining should be lesser or equal to 1625690400`,
		`http response header x-ratelimit-remaining shouldn't be any of [41,49,50]`,    // type conversion isn't supported
		`http response header x-ratelimit-remaining should be any of ["41","49","50"]`, // type conversion isn't supported
	})
}

func doTestV3Api(m *testing.T, id string, opts []string) {
	for _, assertion := range opts {
		name := assertion
		if len(name) > 35 {
			name = name[:32] + "..."
		}
		m.Run(name, func(t *testing.T) {
			doAssertHttpGetProduct(t, id, []byte(fmt.Sprintf(`
		Feature: 
			Scenario:
			Given a HttpRequest
			And the http request headers are:
				| content-type | application/json |
			And the http request method is GET
			And http request URL is https://localhost:8443
			And http request path is /items/`+id+` 
			When the client submits the HttpRequest
			Then %s`, assertion)), nil)
		})
	}
}

func doAssertHttpGetProduct(t *testing.T, id string, def []byte, fm map[string]func() ([]byte, error)) {
	r := readProductFromFile(id)
	m := make(map[string]func() ([]byte, error))
	if fm != nil {
		for k, v := range fm {
			m[k] = v
		}
	}
	m["requests/product.json"] = func() ([]byte, error) {
		return r, nil
	}
	m["schemas/product.json"] = func() ([]byte, error) {
		return getProductJSONSchema(), nil
	}
	fetchSchemaFromServer := func(request *http.Request) (*http.Response, error) {
		f := m["schemas/product.json"]
		binary, err := f()
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: 200,
			Status:     http.StatusText(200),
			Header: map[string][]string{
				"content-type": {"application/json"},
			},
			Body: io.NopCloser(bytes.NewBuffer(binary)),
		}, nil
	}
	ExecV3(t, def, map[string]func(*http.Request) (*http.Response, error){
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
		"GET http://localhost:8080/schemas/product.json":  fetchSchemaFromServer,
		"GET https://localhost:8443/schemas/product.json": fetchSchemaFromServer,
	}, m)
}

func ExecV3(t *testing.T, definition []byte, c map[string]func(*http.Request) (*http.Response, error), fm map[string]func() ([]byte, error)) {
	filesystem := &FnFS{
		fm,
	}
	if fm != nil {
		filesystem.m = fm
	}
	httpClient := &FnHttpClient{
		c,
	}
	ctx := &api.ScenarioDefinitionContext{
		FileSystem: filesystem,
		OnScenarioCreation: func(context *api.DefaultScenarioContext) {
			_ = context.SetValue(ComponentType, "httpClient", httpClient)
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
			ctx.Configure(goDogCtx)
		},
	}
	suite.Run()
}

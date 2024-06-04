package internal

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

//go:embed testdata/features/voucher-api/*
var voucherApiHome embed.FS

func TestScenario_buy_psn_100_uk(t *testing.T) {
	id := "de400a82-c777-4160-b7a0-577a2a0daeef"
	httpClient := FnHttpClient{
		r: map[string]func() ([]byte, error){
			"POST https://localhost:8443/vouchers": func() ([]byte, error) {
				return []byte(fmt.Sprintf(`
				{
					"id":"%s",
					"benefit":"PSN 100 UK",
					"price": {
						"amount": 85,
						"currency": "GBP"
					},
					"has_discount": true
				}`, id)), nil
			},
			"GET https://localhost:8443/vouchers/" + id: func() ([]byte, error) {
				return []byte(fmt.Sprintf(`
				{
					"id":"%s"
				}`, id)), nil
			},
		},
	}
	ctx := &MockContext{
		httpClient:         httpClient,
		resourceHttpClient: httpClient,
		serverURL:          "https://localhost:8443",
		workDirectory:      "testdata/features/voucher-api/",
	}
	content, err := voucherApiHome.ReadFile("testdata/features/voucher-api/feature_with_criteria.feature")
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
					Contents: content,
					Name:     "Buy voucher",
				},
			},
			Output: colors.Colored(os.Stdout),
		},
		ScenarioInitializer: func(scenarioContext *godog.ScenarioContext) {
			_ = Configure(ctx, scenarioContext)
		},
	}.Run()
	assert.Equal(t, 0, status)
}

type FnHttpClient struct {
	r map[string]func() ([]byte, error)
}

func (f FnHttpClient) Do(request *http.Request) (*http.Response, error) {
	key := request.Method + " " + request.URL.String()
	fn, hasValue := f.r[key]
	if !hasValue {
		return &http.Response{
			StatusCode: 404,
			Header:     make(http.Header),
			Status:     http.StatusText(404),
		}, nil
	}
	binary, err := fn()
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Status:     http.StatusText(200),
		Body:       io.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}

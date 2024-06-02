package internal

import (
	"embed"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//go:embed testdata/features/*
var homeDirectory embed.FS

type JUnitOpts struct {
	numberOfScenarios int
	featureName       string
	filename          string
	workDirectory     string
}

func TestHttpRequest_with_single_scenario(t *testing.T) {
	invokeDogAPIGetBreeds(t, JUnitOpts{featureName: "Feature with single scenario", filename: "feature_with_single_scenario.feature", numberOfScenarios: 1})
}

func TestHttpRequest_with_multiple_scenarios(m *testing.T) {
	seq := []string{
		"feature_with_multiple_scenarios.feature",
		"feature_with_multiple_scenarios_with_background.feature",
	}
	for _, filename := range seq {
		m.Run(filename, func(t *testing.T) {
			invokeDogAPIGetBreeds(t, JUnitOpts{
				numberOfScenarios: 3,
				workDirectory:     "testdata/features",
				featureName:       "Feature with multiple scenarios",
				filename:          filename,
			})
		})
	}
}

func TestHttpRequest_with_scenario_with_then_criteria(t *testing.T) {
	invokeDogAPIGetBreed(t, JUnitOpts{featureName: "Criteria API", filename: "feature_with_criteria.feature", numberOfScenarios: 1})
}

func invokeDogAPIGetBreed(t *testing.T, opts JUnitOpts) {
	id := "f9643a80-af1d-422a-9f15-18d466822053"
	serverURL := "https://dogapi.dog/docs/api-v2/breeds/" + id
	httpClient := &SimpleResponseHttpClient{
		statusCode: 200,
		fileURI:    "features/GetBreed.json",
		headers: map[string]string{
			"content-type":          "application/json",
			"x-ratelimit-remaining": "2",
			"x-ratelimit-limit":     "100",
			"x-ratelimit-reset":     "1625074801",
		},
	}
	content, err := homeDirectory.ReadFile("testdata/features/" + opts.filename)
	if err != nil {
		t.Fatal(err)
	}
	goDogOpts := godog.Options{
		TestingT: t,
		Format:   "pretty",
		Paths:    []string{},
		FeatureContents: []godog.Feature{
			{
				Contents: content,
				Name:     opts.featureName,
			},
		},
		Output: colors.Colored(os.Stdout),
	}
	workDirectory := opts.workDirectory
	if workDirectory == "" {
		if w, prob := os.Getwd(); prob == nil {
			workDirectory = w
		} else {
			t.Fatal(err)
		}
	}
	status := godog.TestSuite{
		TestSuiteInitializer: nil,
		Options:              &goDogOpts,
		ScenarioInitializer: New(&MockContext{
			serverURL:     serverURL,
			httpClient:    httpClient,
			workDirectory: workDirectory,
			resourceHttpClient: &EmbeddedResourceHttpClient{
				statusCode: 200,
				directory:  "features",
				fs:         homeDirectory,
			},
		}),
	}.Run()
	assert.Equal(t, 0, status)
	assert.Equal(t, len(httpClient.Requests), opts.numberOfScenarios)
	for i := range httpClient.Requests {
		request := httpClient.Requests[i]
		// Assert Body
		assert.Nil(t, request.Body)
		// Assert Headers
		assert.Len(t, request.Header, 1)
		assert.Equal(t, request.Header.Get("content-type"), "application/json")
		// Assert method and URL
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, serverURL+"/breeds/"+id, request.URL.String())
	}
}

func invokeDogAPIGetBreeds(t *testing.T, opts JUnitOpts) {
	serverURL := "https://dogapi.dog/docs/api-v2/breeds"
	httpClient := &SimpleResponseHttpClient{
		statusCode: 200,
		fileURI:    "features/GetBreeds.json",
		headers:    map[string]string{"content-type": "application/json"},
	}
	content, err := homeDirectory.ReadFile("testdata/features/" + opts.filename)
	if err != nil {
		t.Fatal(err)
	}
	goDogOpts := godog.Options{
		TestingT: t,
		Format:   "pretty",
		Paths:    []string{},
		FeatureContents: []godog.Feature{
			{
				Contents: content,
				Name:     opts.featureName,
			},
		},
		Output: colors.Colored(os.Stdout),
	}
	workDirectory := opts.workDirectory
	if workDirectory == "" {
		if w, prob := os.Getwd(); prob == nil {
			workDirectory = w
		} else {
			t.Fatal(err)
		}
	}
	status := godog.TestSuite{
		TestSuiteInitializer: nil,
		Options:              &goDogOpts,
		ScenarioInitializer: New(&MockContext{
			serverURL:     serverURL,
			httpClient:    httpClient,
			workDirectory: workDirectory,
			resourceHttpClient: &EmbeddedResourceHttpClient{
				statusCode: 200,
				directory:  "features",
				fs:         homeDirectory,
			},
		}),
	}.Run()
	assert.Equal(t, 0, status)
	assert.Equal(t, len(httpClient.Requests), opts.numberOfScenarios)
	for i := range httpClient.Requests {
		request := httpClient.Requests[i]
		// Assert Body
		assert.Nil(t, request.Body)
		// Assert Headers
		assert.Len(t, request.Header, 1)
		assert.Equal(t, request.Header.Get("content-type"), "application/json")
		// Assert method and URL
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, serverURL+"/breeds", request.URL.String())
	}
}

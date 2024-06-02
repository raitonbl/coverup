package internal

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/raitonbl/coverup/pkg"
	"io"
	"net/http"
)

type MockContext struct {
	serverURL          string
	workDirectory      string
	httpClient         pkg.HttpClient
	resourceHttpClient pkg.HttpClient
}

func (m *MockContext) GetServerURL() string {
	return m.serverURL
}

func (m *MockContext) GetWorkDirectory() string {
	return m.workDirectory
}

func (m *MockContext) GetHttpClient() pkg.HttpClient {
	return m.httpClient
}

func (m *MockContext) GetResourcesHttpClient() pkg.HttpClient {
	return m.resourceHttpClient
}

type EmbeddedResourceHttpClient struct {
	statusCode int
	directory  string
	fs         embed.FS
	headers    map[string]string
}

func (instance *EmbeddedResourceHttpClient) Do(req *http.Request) (*http.Response, error) {
	f := "testdata/"
	if instance.directory != "" {
		f += instance.directory
	}
	f += req.URL.Path
	fmt.Println("FILE:", f)
	content, err := instance.fs.ReadFile(f)
	if err != nil {
		return nil, err
	}
	// Create a response object
	response := &http.Response{
		Header:     make(http.Header),
		StatusCode: instance.statusCode,
		Status:     http.StatusText(instance.statusCode),
		Body:       io.NopCloser(bytes.NewBuffer(content)),
	}
	// Set the headers
	if instance.headers != nil {
		for key, value := range instance.headers {
			response.Header.Set(key, value)
		}
	}
	return response, nil
}

type SimpleResponseHttpClient struct {
	statusCode int
	err        error
	content    []byte
	fileURI    string
	Requests   []http.Request
	headers    map[string]string
	Filter     func(r *http.Request) bool
}

func (instance *SimpleResponseHttpClient) Do(req *http.Request) (*http.Response, error) {
	if instance.Requests == nil {
		instance.Requests = make([]http.Request, 0)
	}
	if instance.Filter == nil || instance.Filter(req) {
		instance.Requests = append(instance.Requests, *req)
	}
	if instance.err != nil {
		return nil, instance.err
	}
	content := instance.content
	if content == nil && instance.fileURI != "" {
		if c, err := dogsApiHomeDirectory.ReadFile("testdata/" + instance.fileURI); err == nil {
			content = c
		} else {
			return nil, err
		}
	}
	// Create a response object
	response := &http.Response{
		Header:     make(http.Header),
		StatusCode: instance.statusCode,
		Status:     http.StatusText(instance.statusCode),
		Body:       io.NopCloser(bytes.NewBuffer(content)),
	}
	// Set the headers
	if instance.headers != nil {
		for key, value := range instance.headers {
			response.Header.Set(key, value)
		}
	}
	return response, nil
}

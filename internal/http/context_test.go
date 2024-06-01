package internal

import (
	"bytes"
	"io"
	"net/http"
)

type MockContext struct {
	serverURL     string
	workDirectory string
	httpClient    HttpClient
}

func (m *MockContext) GetServerURL() string {
	return m.serverURL
}

func (m *MockContext) GetWorkDirectory() string {
	return m.workDirectory
}

func (m *MockContext) GetHttpClient() HttpClient {
	return m.httpClient
}

type SimpleResponseHttpClient struct {
	err        error
	content    []byte
	fileURI    string
	statusCode int
	headers    map[string]string
	Requests   []http.Request
}

func (instance *SimpleResponseHttpClient) Do(req *http.Request) (*http.Response, error) {
	if instance.Requests == nil {
		instance.Requests = make([]http.Request, 0)
	}
	instance.Requests = append(instance.Requests, *req)
	if instance.err != nil {
		return nil, instance.err
	}
	content := instance.content
	if content == nil && instance.fileURI != "" {
		if c, err := homeDirectory.ReadFile("testdata/" + instance.fileURI); err == nil {
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

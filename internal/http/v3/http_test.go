package v3

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	array := []map[string]interface{}{
		{"id": "05f584e5-41d9-448d-ad39-321a39badd92"},
		{"id": "fc6eeb12-f3ce-4dec-92dc-c6f3030f82bf"},
		{"id": "fb1c9978-ab5f-4bb2-bc7f-163e245656aa"},
	}
	binary, err := json.Marshal(array)
	if err != nil {
		t.Error(err)
	}
	Exec(t, []byte(`
		Feature: 
			Scenario:
			Given a HttpRequest
			And the headers:
				| content-type | application/json |
			And Method is GET
			And server url is https://localhost:8443 
			And Path is /items
			When the client submits the HttpRequest
	`), map[string]func(*http.Request) (*http.Response, error){
		"GET https://localhost:8443/items": func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Header: map[string][]string{
					"content-type": {"application/json"},
				},
				Body: io.NopCloser(bytes.NewBuffer(binary)),
			}, nil
		},
	}, make(map[string][]byte))
}

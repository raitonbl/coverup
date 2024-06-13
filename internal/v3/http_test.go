package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestHttp_AssertResponseStatusCode(t *testing.T) {
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
			Then response status code is 200
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
	}, nil)
}

func TestHttpContext_AssertResponseContentType(t *testing.T) {
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
			Then response content-type is "application/json"
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
	}, nil)
}

func TestHttpContext_AssertResponseIsValidAgainstSchema(t *testing.T) {
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
			Then response content-type is "application/json"
			And response body respects schema file://schema.json
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
	}, map[string]func() ([]byte, error){
		"schema.json": func() ([]byte, error) {
			return []byte(`
				{
				  "$schema": "http://json-schema.org/draft-07/schema#",
				  "type": "array",
				  "items": {
					"type": "object",
					"properties": {
					  "id": {
						"type": "string",
						"pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
					  }
					},
					"required": ["id"],
					"additionalProperties": false
				  }
				}
			`), nil
		},
	})
}

func TestHttpContext_AssertSimpleAttribute(m *testing.T) {
	id := "27258303-9ebc-4b84-a17e-f886161ab2f5"
	opts := []string{
		`response body $.id is "` + id + `"`,
		`response body $.id isn't "RaitonBL"`,
		`response body $.summary contains "SSD"`,
		`response body $.summary ignoring case contains "ssd"`,
		`response body $.name starts with "Seagate"`,
		`response body $.name ignoring case starts with "SeaGate"`,
		`response body $.name ends with "MB/s"`,
		`response body $.name ignoring case ends with "mb/s"`,
		`response body $.warranty.amount is 2`,
		`response body $.warranty.amount isn't 3`,
		`response body $.in_promotion is false`,
		`response body $.in_promotion isn't true`,
		`response body $.name matches pattern "^Seagate"`,
		`response body $.name doesn't match pattern "^X"`,
		`response body respects schema file://schemas/product.json`,
		`response body respects schema http://localhost:8080/schemas/product.json`,
		`response body respects schema https://localhost:8443/schemas/product.json`,
		`response body is:
		"""
			` + string(readProductFromFile(id)) + `
		"""`,
		`response body is file://requests/product.json`,
	}
	for _, assertion := range opts {
		name := assertion
		if len(name) > 35 {
			name = name[:32] + "..."
		}
		m.Run(name, func(t *testing.T) {
			assertHttpGetProduct(t, id, []byte(fmt.Sprintf(`
		Feature: 
			Scenario:
			Given a HttpRequest
			And the headers:
				| content-type | application/json |
			And Method is GET
			And server url is https://localhost:8443
			And Path is /items/`+id+` 
			When the client submits the HttpRequest
			Then response content-type is "application/json"
			And %s
	`, assertion)), nil)
		})
	}

}

func assertHttpGetProduct(t *testing.T, id string, def []byte, fm map[string]func() ([]byte, error)) {
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
	Exec(t, def, map[string]func(*http.Request) (*http.Response, error){
		fmt.Sprintf("GET https://localhost:8443/items/%s", id): func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Header: map[string][]string{
					"content-type": {"application/json"},
				},
				Body: io.NopCloser(bytes.NewBuffer(r)),
			}, nil
		},
		"GET http://localhost:8080/schemas/product.json":  fetchSchemaFromServer,
		"GET https://localhost:8443/schemas/product.json": fetchSchemaFromServer,
	}, m)
}

func readProductFromFile(id string) []byte {
	return []byte(`
	{
	  "id": "` + id + `",
	  "name": "Seagate One Touch SSD 1TB External SSD Portable – Black, speeds up to 1030MB/s",
	  "summary": "An external SSD device that guarantees 1TB storage",
	  "image": "/photos/27258303-9ebc-4b84-a17e-f886161ab2f5",
	  "in_promotion": false,
	  "offer_created_at": "2022-06-06T12:34:56Z",
	  "offer_expires_at": "2024-12-31T23:59:56Z",
	  "about": [
		"One Touch SSD is a mini USB 3.0 SSD featuring a lightweight, textile design for busy days and bustling commutes.",
		"High-speed, portable solid state drive perfect for streaming stored videos directly to laptop, scrolling seamlessly through photos, and backing up content on the go. ",
		"Enjoy long-term peace of mind with the included three-year limited warranty and Rescue Data Recovery Services. "
	  ],
	  "tags": [
		{
		  "id": "ebbb5082-58f4-4ea4-9840-e02cc86501de",
		  "name": "IT"
		},
		{
		  "id": "0c61ba6a-baea-4316-b2c2-e847253d029b",
		  "name": "HDD"
		}
	  ],
	  "warranty": {
		"amount": 2,
		"unit": "years"
	  },
	  "price": {
		"amount": 200,
		"currency": "USD"
	  },
	  "characteristics": {
		"capacity": {
		  "amount": 1,
		  "unit": "TB"
		},
		"hard_disk_interface": "USB-C",
		"connectivity_technology": "USB",
		"brand": "Seagate",
		"special_feature": "Portable",
		"hard_disk_form_factor": {
		  "amount": 2.5,
		  "unit": "inch"
		},
		"hard_disk_description": "SSD",
		"color": "BLACK",
		"installation_type": "EXTERNAL"
	  }
	}`)
}

func getProductJSONSchema() []byte {
	return []byte(`
			{
			  "$schema": "http://json-schema.org/draft-07/schema#",
			  "type": "object",
			  "properties": {
				"id": {
				  "type": "string"
				},
				"name": {
				  "type": "string"
				},
				"summary": {
				  "type": "string"
				},
				"image": {
				  "type": "string"
				},
				"in_promotion": {
				  "type": "boolean"
				},
				"offer_created_at": {
				  "type": "string",
				  "format": "date-time"
				},
				"offer_expires_at": {
				  "type": "string",
				  "format": "date-time"
				},
				"about": {
				  "type": "array",
				  "items": {
					"type": "string"
				  }
				},
				"tags": {
				  "type": "array",
				  "items": {
					"type": "object",
					"properties": {
					  "id": {
						"type": "string"
					  },
					  "name": {
						"type": "string"
					  }
					},
					"required": ["id", "name"]
				  }
				},
				"warranty": {
				  "type": "object",
				  "properties": {
					"amount": {
					  "type": "integer"
					},
					"unit": {
					  "type": "string"
					}
				  },
				  "required": ["amount", "unit"]
				},
				"price": {
				  "type": "object",
				  "properties": {
					"amount": {
					  "type": "number"
					},
					"currency": {
					  "type": "string"
					}
				  },
				  "required": ["amount", "currency"]
				},
				"characteristics": {
				  "type": "object",
				  "properties": {
					"capacity": {
					  "type": "object",
					  "properties": {
						"amount": {
						  "type": "integer"
						},
						"unit": {
						  "type": "string"
						}
					  },
					  "required": ["amount", "unit"]
					},
					"hard_disk_interface": {
					  "type": "string"
					},
					"connectivity_technology": {
					  "type": "string"
					},
					"brand": {
					  "type": "string"
					},
					"special_feature": {
					  "type": "string"
					},
					"hard_disk_form_factor": {
					  "type": "object",
					  "properties": {
						"amount": {
						  "type": "number"
						},
						"unit": {
						  "type": "string"
						}
					  },
					  "required": ["amount", "unit"]
					},
					"hard_disk_description": {
					  "type": "string"
					},
					"color": {
					  "type": "string"
					},
					"installation_type": {
					  "type": "string"
					}
				  },
				  "required": [
					"capacity",
					"hard_disk_interface",
					"connectivity_technology",
					"brand",
					"special_feature",
					"hard_disk_form_factor",
					"hard_disk_description",
					"color",
					"installation_type"
				  ]
				}
			  },
			  "required": [
				"id",
				"name",
				"summary",
				"image",
				"in_promotion",
				"offer_created_at",
				"offer_expires_at",
				"about",
				"tags",
				"warranty",
				"price",
				"characteristics"
			  ]
			}
		`)
}

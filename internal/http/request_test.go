package internal

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func TestFn(t *testing.T) {
	fmt.Println(Get("Operation GET /vouchers"))
	fmt.Println(Get("Operation GET /vouchers/{{HttpRequest.SendVoucherRequest.Response.Body.id}}"))
	fmt.Println(Get("Operation GET /vouchers/{{HttpRequest.SendVoucherRequest.Response.Body.id}}/{{HttpRequest.SendVoucherRequest.Response.Body.name}}"))
}

func Get(p string) (any, error) {
	var valueRegexp = regexp.MustCompile(`\{\{([^.]+)\.([^.]+)\.(.+?)\}\}`)

	var prob error
	result := valueRegexp.ReplaceAllStringFunc(p, func(currentValue string) string {
		s := valueRegexp.FindStringSubmatch(currentValue)
		if len(s) == 0 {
			return currentValue
		}
		if len(s) != 4 {
			prob = errors.New("cannot resolve " + p)
			return ""
		}
		// Extract the captured groups
		componentType := s[1]
		componentId := s[2]
		anything := s[3]

		// Print the captured groups

		// Create the replacement string
		replacement := componentType + " > " + componentId + " > " + anything

		return replacement
	})
	if prob != nil {
		return nil, prob
	}
	return result, nil
}

func TestResolveResponse(t *testing.T) {
	name := "RaitonBL"
	id := "bd3cbbad-599b-4676-9fc5-11b639656f1d"
	response := Response{
		statusCode: 200,
		body: []byte(fmt.Sprintf(`
		{
			"id":"%s",
			"name":"%s"
		}
		`, id, name)),
		pathCache: make(map[string]any),
		headers: map[string]string{
			"content-type": "application/json",
		},
	}
	request := Request{
		method:    "GET",
		uri:       "/info",
		serverURL: "https://localhost:9443",
		headers:   map[string]string{"content-type": "application/json"},
		response:  &response,
	}
	valueOf, err := request.GetPathValue("Response.Body.id")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(valueOf)
}

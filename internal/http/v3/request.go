package v3

type HttpRequest struct {
	form      *Form
	path      string
	method    string
	serverURL string
	body      []byte
	headers   map[string]string
	response  *HttpResponse
}

func (instance *HttpRequest) GetPathValue(x string) (any, error) {
	panic("implement me")
}

type Form struct {
	encType    string
	attributes map[string]any
}

func (instance *Form) GetPathValue(x string) (any, error) {
	panic("implement me")
}

type HttpResponse struct {
	body       []byte
	headers    map[string]string
	statusCode int
	pathCache  map[string]any
}

func (h HttpResponse) GetPathValue(x string) (any, error) {
	panic("implement me")
}

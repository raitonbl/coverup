package v3

type HttpRequest struct {
	form      *Form
	path      string
	method    string
	serverURL string
	body      []byte
	headers   map[string]string
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

package v3

type HttpRequest struct {
	form      *Form
	path      string
	method    string
	serverURL string
	body      []byte
	headers   map[string]string
}

type Form struct {
	encType    string
	attributes map[string][]byte
}

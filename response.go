package heimdall

// Response encapsulates details of a http response
type Response struct {
	body       []byte
	statusCode int
}

// StatusCode returns status code of a http request
func (hr Response) StatusCode() int {
	return hr.statusCode
}

// Body returns body in bytes of a http request
func (hr Response) Body() []byte {
	return hr.body
}

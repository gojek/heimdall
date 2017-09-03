package heimdall

type HeimdallResponse struct {
	body       []byte
	statusCode int
}

func (hr HeimdallResponse) StatusCode() int {
	return hr.statusCode
}

func (hr HeimdallResponse) Body() []byte {
	return hr.body
}

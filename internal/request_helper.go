package internal

import (
	"bytes"
	"io"
	"net/http"
)

type RequestGetBody func() (io.ReadCloser, error)

func SetRequestGetBody(r *http.Request) error {
	if r == nil || r.Body == nil || r.GetBody != nil { // skip if body not present or GetBody is already present
		return nil
	}

	if r.Body == http.NoBody { // optimized handling for NoBody cases
		r.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		return nil
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body = newBytesBody(buf)
	r.GetBody = func() (io.ReadCloser, error) {
		return newBytesBody(buf), nil
	}

	return nil
}

func CloneRequest(request *http.Request, getReqBody RequestGetBody) (*http.Request, error) {
	if request == nil || getReqBody == nil {
		return request, nil
	}

	body, err := getReqBody()
	if err != nil {
		return nil, err
	}
	request = request.WithContext(request.Context()) // shallow clone instead of deep clone from http.Request.Clone
	request.Body = body
	return request, nil
}

func newBytesBody(buf []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(buf))
}

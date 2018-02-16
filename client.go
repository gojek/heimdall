package heimdall

import (
	"io"
	"net/http"
)

// Client Is a generic HTTP client interface
type Client interface {
	Get(url string, headers http.Header) (Response, error)
	Post(url string, body io.Reader, headers http.Header) (Response, error)
	Put(url string, body io.Reader, headers http.Header) (Response, error)
	Patch(url string, body io.Reader, headers http.Header) (Response, error)
	Delete(url string, headers http.Header) (Response, error)
	Do(req *http.Request) (*http.Response, error)

	SetRetryCount(count int)
	SetRetrier(retrier Retriable)
}

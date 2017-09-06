package heimdall

import "io"

// Client Is a generic HTTP client interface
type Client interface {
	Get(url string) (Response, error)
	Post(url string, body io.Reader) (Response, error)
	Put(url string, body io.Reader) (Response, error)
	Patch(url string, body io.Reader) (Response, error)
	Delete(url string) (Response, error)

	SetRetryCount(count int)
	SetRetrier(retrier Retriable)
}

package heimdall

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client Is a generic client interface
type Client interface {
	Get(url string) (Response, error)
}

type httpClient struct {
	client *http.Client
}

// NewHTTPClient returns a new instance of HTTPClient
func NewHTTPClient(config Config) Client {
	timeout := config.timeoutInSeconds

	httpTimeout := time.Duration(timeout) * time.Second
	return &httpClient{
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// Get makes a HTTP get request to provided URL
func (c *httpClient) Get(url string) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	return c.do(request)
}

func (c *httpClient) do(request *http.Request) (Response, error) {
	hr := Response{}
	var err error

	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		return hr, err
	}

	defer response.Body.Close()

	var responseBytes []byte
	if response.Body != nil {
		responseBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return hr, err
		}
	}

	hr.body = responseBytes
	hr.statusCode = response.StatusCode

	return hr, err
}

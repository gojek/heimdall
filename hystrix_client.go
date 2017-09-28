package heimdall

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/pkg/errors"
)

const defaultHystrixRetryCount int = 0

type hystrixHTTPClient struct {
	client *http.Client

	hystrixCommandName string

	retryCount int
	retrier    Retriable
}

// NewHystrixHTTPClient returns a new instance of HystrixHTTPClient
func NewHystrixHTTPClient(httpClient *http.Client, hystrixConfig HystrixConfig) Client {
	hystrix.ConfigureCommand(hystrixConfig.commandName, hystrixConfig.commandConfig)

	return &hystrixHTTPClient{
		client: httpClient,

		retryCount:         defaultHystrixRetryCount,
		retrier:            NewNoRetrier(),
		hystrixCommandName: hystrixConfig.commandName,
	}
}

// SetRetryCount sets the retry count for the hystrixHTTPClient
func (hhc *hystrixHTTPClient) SetRetryCount(count int) {
	hhc.retryCount = count
}

// SetRetrier sets the strategy for retrying
func (hhc *hystrixHTTPClient) SetRetrier(retrier Retriable) {
	hhc.retrier = retrier
}

// Get makes a HTTP GET request to provided URL
func (hhc *hystrixHTTPClient) Get(url string, headers http.Header) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return hhc.do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Post(url string, body io.Reader, headers http.Header) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return hhc.do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Put(url string, body io.Reader, headers http.Header) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return hhc.do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Patch(url string, body io.Reader, headers http.Header) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return hhc.do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (hhc *hystrixHTTPClient) Delete(url string, headers http.Header) (Response, error) {
	response := Response{}

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return hhc.do(request)
}

func (hhc *hystrixHTTPClient) do(request *http.Request) (Response, error) {
	hr := Response{}

	request.Close = true

	var err error
	for i := 0; i <= hhc.retryCount; i++ {

		err = hystrix.Do(hhc.hystrixCommandName, func() error {
			response, err := hhc.client.Do(request)
			if err != nil {
				return err
			}

			if response.Body != nil {
				hr.body, err = ioutil.ReadAll(response.Body)
				if err != nil {
					return err
				}
			}

			response.Body.Close()

			hr.statusCode = response.StatusCode

			if response.StatusCode >= http.StatusInternalServerError {
				return fmt.Errorf("Server is down: returned status code: %d", response.StatusCode)
			}

			return nil
		}, func(err error) error {
			return err
		})

		if err != nil {
			backoffTime := hhc.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		break
	}

	return hr, err
}

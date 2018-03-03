package examples

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/pkg/errors"
)

const (
	baseURL = "http://localhost:9090"
)

func httpClientUsage() error {
	timeout := 100 * time.Millisecond

	httpClient := heimdall.NewHTTPClient(timeout)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	httpClient.SetRetryCount(2)
	httpClient.SetRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10)))

	response, err := httpClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

func hystrixClientUsage() error {
	timeout := 100 * time.Millisecond

	hystrixConfig := heimdall.NewHystrixConfig("MyCommand", heimdall.HystrixCommandConfig{
		Timeout:                1100,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  25,
		SleepWindow:            10,
		RequestVolumeThreshold: 10,
	})

	hystrixClient := heimdall.NewHystrixHTTPClient(timeout, hystrixConfig)
	headers := http.Header{}
	response, err := hystrixClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

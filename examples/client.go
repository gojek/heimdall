package examples

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"github.com/gojektech/heimdall/hystrix"
	"github.com/pkg/errors"
)

const (
	baseURL = "http://localhost:9090"
)

func httpClientUsage() error {
	timeout := 100 * time.Millisecond

	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetryCount(2),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10, 5))),
	)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

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
	hystrixClient := hystrix.NewClient(
		hystrix.WithHTTPTimeout(timeout),
		hystrix.WithCommandName("MyCommand"),
		hystrix.WithHystrixTimeout(1100),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(25),
		hystrix.WithSleepWindow(10),
		hystrix.WithRequestVolumeThreshold(10),
	)
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

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth("username", "passwd")
	return c.client.Do(request)
}

func customHTTPClientUsage() error {
	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(0*time.Millisecond),
		httpclient.WithHTTPClient(&myHTTPClient{
			// replace with custom HTTP client
			client: http.Client{Timeout: 25 * time.Millisecond},
		}),
		httpclient.WithRetryCount(2),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10, 5))),
	)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

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

func customHystrixClientUsage() error {
	timeout := 0 * time.Millisecond

	hystrixClient := hystrix.NewClient(
		hystrix.WithHTTPTimeout(timeout),
		hystrix.WithCommandName("MyCommand"),
		hystrix.WithHystrixTimeout(1100),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(25),
		hystrix.WithSleepWindow(10),
		hystrix.WithRequestVolumeThreshold(10),
		hystrix.WithHTTPClient(&myHTTPClient{
			// replace with custom HTTP client
			client: http.Client{Timeout: 25 * time.Millisecond},
		}),
	)

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

func httpClientWithExponentialBackoff() error {
	initalTimeout := 2 * time.Millisecond
	maxTimeout := 9 * time.Millisecond
	exponentFactor := 2.0
	maximumJitterInterval := 2 * time.Millisecond

	backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, exponentFactor, maximumJitterInterval)

	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(3),
	)

	headers := http.Header{}
	response, err := client.Get(baseURL, headers)
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

func httpClientWithConstantBackoff() error {
	backOffInterval := 10 * time.Second
	maximumJitterInterval := 2 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backOffInterval, maximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(3),
	)

	headers := http.Header{}
	response, err := client.Get(baseURL, headers)
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

type linearBackoff struct {
	backoffInterval int
}

func (lb *linearBackoff) Next(retry int) time.Duration {
	if retry <= 0 {
		return 0 * time.Millisecond
	}
	return time.Duration(retry*lb.backoffInterval) * time.Millisecond
}

func httpClientWithCustomBackoff() error {
	backoff := &linearBackoff{100}
	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(3),
	)

	headers := http.Header{}
	response, err := client.Get(baseURL, headers)
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

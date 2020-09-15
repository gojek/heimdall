package hystrix

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Set("foo", "bar")
	return c.client.Do(request)
}

func TestHystrixHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(50*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(20),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")

	response, err := client.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHystrixHTTPClientGetSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Get(server.URL, headers)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPostSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Post(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientDeleteSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Delete(server.URL, headers)
	require.NoError(t, err, "should not have failed to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPutSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Put(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPatchSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	requestBody := bytes.NewReader([]byte(requestBodyString))

	response, err := client.Patch(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientRetriesGetOnFailure(t *testing.T) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	response, err := client.Get("url_doesnt_exist", http.Header{})

	assert.Contains(t, err.Error(), "unsupported protocol scheme")
	assert.Nil(t, response)
}

func TestHystrixHTTPClientRetriesGetOnFailure5xx(t *testing.T) {
	count := 0
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name_5xx"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
		count = count + 1
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, 4, count)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))
}

func BenchmarkHystrixHTTPClientRetriesGetOnFailure(b *testing.B) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Get(server.URL, http.Header{})
	}
}

func TestHystrixHTTPClientDontRetryWhenContextIsCancelled(t *testing.T) {
	t.Skip("Skip: concurrency issues... to be fixed in hystrix")

	count := 0
	noOfRetries := 3
	// Set a huge backoffInterval that we won't have to wait anyway
	backoffInterval := 1 * time.Hour
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	ctx, cancel := context.WithCancel(context.Background())

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		// Cancel the context after the first call
		cancel()

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response": "something went wrong"}`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	response, err := client.Do(req.WithContext(ctx))
	require.Error(t, err, "should have failed to make request")
	require.Nil(t, response)

	assert.Equal(t, 1, count)
}

func TestHystrixHTTPClientRetriesPostOnFailure(t *testing.T) {
	count := 0
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(50*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(20),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
		count = count + 1
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Post(server.URL, strings.NewReader("a=1&b=2"), http.Header{})
	require.NoError(t, err)

	assert.Equal(t, 4, count)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.JSONEq(t, `{ "response": "something went wrong" }`, respBody(t, response))
}

func BenchmarkHystrixHTTPClientRetriesPostOnFailure(b *testing.B) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Post(server.URL, strings.NewReader("a=1&b=2"), http.Header{})
	}
}

func TestHystrixHTTPClientReturnsFallbackFailureWithoutFallBackFunction(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Equal(t, err.Error(), "hystrix: circuit open")
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnAnError(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			// do something in the fallback function
			return err
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, strings.Contains(err.Error(), "fallback failed"))
}

func TestFallBackFunctionIsCalledWithHystrixHTTPClient(t *testing.T) {
	called := false

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			called = true
			return err
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, called)
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnsNil(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			// do something in the fallback function
			return nil
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Nil(t, err)
}

func TestCustomHystrixHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_new_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithHTTPClient(&myHTTPClient{
			client: http.Client{Timeout: 25 * time.Millisecond}}),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("foo"), "bar")
		assert.NotEqual(t, r.Header.Get("foo"), "baz")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, err := client.Do(req)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer func() {
			_ = response.Body.Close()
		}()
	}

	respBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")

	return string(respBody)
}

func TestDurationToInt(t *testing.T) {
	t.Run("1sec should return 1 when unit is second", func(t *testing.T) {
		timeout := 1 * time.Second
		timeoutInSec := durationToInt(timeout, time.Second)

		assert.Equal(t, 1, timeoutInSec)
	})

	t.Run("30sec should return 30000 when unit is millisecond", func(t *testing.T) {
		timeout := 30 * time.Second
		timeoutInMs := durationToInt(timeout, time.Millisecond)

		assert.Equal(t, 30000, timeoutInMs)
	})
}

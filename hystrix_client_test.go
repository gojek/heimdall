package heimdall

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afex/hystrix-go/hystrix"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHystrixHTTPClientDoSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Accept-Language"), "en")

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
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

func TestHystrixHTTPClientRetriesOnFailure(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	count := 0

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count = count + 1
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	client.SetRetryCount(3)
	client.SetRetrier(NewRetrier(NewConstantBackoff(1)))

	response, err := client.Get(server.URL, http.Header{})
	require.Error(t, err)

	assert.Equal(t, 4, count)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))
}

func TestHystrixHTTPClientReturnsFallbackFailureWithoutFallBackFunction(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
	})

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Equal(t, err.Error(), "hystrix: circuit open")
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnAnError(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
		fallbackFn: func(err error) error {
			// do something in the fallback function
			return err
		},
	})

	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, strings.Contains(err.Error(), "fallback failed"))
}

func TestFallBackFunctionIsCalledWithHystrixHTTPClient(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	called := false
	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
		fallbackFn: func(err error) error {
			called = true
			return err
		},
	})
	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, called)
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnsNil(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, HystrixConfig{
		commandName:   "some_command_name",
		commandConfig: hystrixCommandConfig,
		fallbackFn: func(err error) error {
			// do something in the fallback function
			return nil
		},
	})

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Nil(t, err)
}

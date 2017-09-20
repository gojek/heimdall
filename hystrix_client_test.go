package heimdall

import (
	"bytes"
	"github.com/afex/hystrix-go/hystrix"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

func TestHystrixHTTPClientGetSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}

	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Not a GET request")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHystrixHTTPClientPostSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Not a POST request")
		}

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		if string(rBody) != requestBodyString {
			t.Errorf("POST request has wrong request body")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	response, err := client.Post(server.URL, requestBody)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHystrixHTTPClientDeleteSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Not a DELETE request")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Delete(server.URL)
	require.NoError(t, err, "should not have failed to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHystrixHTTPClientPutSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Not a PUT request")
		}

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		if string(rBody) != requestBodyString {
			t.Errorf("PUT request has wrong request body")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	response, err := client.Put(server.URL, requestBody)
	require.NoError(t, err, "should not have failed to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHystrixHTTPClientPatchSuccess(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Not a PATCH request")
		}

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		if string(rBody) != requestBodyString {
			t.Errorf("PATCH request has wrong request body")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	response, err := client.Patch(server.URL, requestBody)
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHystrixHTTPClientRetriesOnFailure(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

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

	response, err := client.Get(server.URL)
	require.Error(t, err)

	assert.Equal(t, 4, count)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"something went wrong\" }", string(response.Body()))
}

func TestHystrixHTTPClientReturnsFallbackFailure(t *testing.T) {
	hystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  10,
		SleepWindow:            100,
		RequestVolumeThreshold: 10,
	}
	client := NewHystrixHTTPClient(10, hystrixCommandConfig)

	_, err := client.Get("http://localhost")
	assert.True(t, strings.Contains(err.Error(), "fallback failed"))
}

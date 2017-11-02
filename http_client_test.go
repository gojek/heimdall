package heimdall

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientGetSuccess(t *testing.T) {
	client := NewHTTPClient(10)

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

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHTTPClientPostSuccess(t *testing.T) {
	client := NewHTTPClient(10)

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

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHTTPClientDeleteSuccess(t *testing.T) {
	client := NewHTTPClient(10)

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

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHTTPClientPutSuccess(t *testing.T) {
	client := NewHTTPClient(10)

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

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHTTPClientPatchSuccess(t *testing.T) {
	client := NewHTTPClient(10)

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

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Patch(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, "{ \"response\": \"ok\" }", string(response.Body()))
}

func TestHTTPClientGetRetriesOnFailure(t *testing.T) {
	client := NewHTTPClient(10)

	count := 0

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	noOfRetries := 3
	noOfCalls := noOfRetries + 1

	client.SetRetryCount(noOfRetries)
	client.SetRetrier(NewRetrier(NewConstantBackoff(1)))

	response, err := client.Get(server.URL, http.Header{})

	require.Equal(t, http.StatusInternalServerError, response.StatusCode())
	require.Equal(t, "{ \"response\": \"something went wrong\" }", string(response.Body()))

	assert.Equal(t, noOfCalls, count)
	assert.Error(t, err)
}

func TestHTTPClientGetReturnsAllErrorsIfRetriesFail(t *testing.T) {
	client := NewHTTPClient(10)

	count := 0

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	noOfRetries := 2
	client.SetRetryCount(noOfRetries)
	client.SetRetrier(NewRetrier(NewConstantBackoff(1)))

	response, err := client.Get(server.URL, http.Header{})

	require.Equal(t, noOfRetries+1, count)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode())
	require.Equal(t, "{ \"response\": \"something went wrong\" }", string(response.Body()))

	assert.Equal(t, "server error: 500, server error: 500, server error: 500", err.Error())
}

func TestHTTPClientGetReturnsNoErrorsIfRetrySucceeds(t *testing.T) {
	client := NewHTTPClient(10)

	count := 0
	countWhenCallSucceeds := 2

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if count == countWhenCallSucceeds {
			w.Write([]byte(`{ "response": "success" }`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "response": "something went wrong" }`))
		}
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	client.SetRetryCount(3)
	client.SetRetrier(NewRetrier(NewConstantBackoff(1)))

	response, err := client.Get(server.URL, http.Header{})

	require.Equal(t, countWhenCallSucceeds+1, count)
	require.Equal(t, http.StatusOK, response.StatusCode())
	require.Equal(t, "{ \"response\": \"success\" }", string(response.Body()))

	assert.NoError(t, err)
}

func TestHTTPClientGetReturnsErrorOnClientCallFailure(t *testing.T) {
	client := NewHTTPClient(10)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	server.URL = "" // Invalid URL to simulate client.Do error
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})

	require.NotEqual(t, http.StatusOK, response.StatusCode())

	assert.Equal(t, "Get : unsupported protocol scheme \"\"", err.Error())
}

func TestHTTPClientGetReturnsErrorOnParseResponseFailure(t *testing.T) {
	client := NewHTTPClient(10)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate unexpected EOF with a response longer than Content-Length
		w.Header().Set("Content-Length", "3")
		w.Write([]byte("aasd"))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})

	require.NotEqual(t, http.StatusOK, response.StatusCode())

	assert.Equal(t, "unexpected EOF", err.Error())
}

func TestHTTPClientGetReturnsErrorOn5xxFailure(t *testing.T) {
	client := NewHTTPClient(10)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})

	require.Equal(t, http.StatusInternalServerError, response.StatusCode())

	assert.Equal(t, "server error: 500", err.Error())
}

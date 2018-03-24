package heimdall

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientWithContextSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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
	response, err := client.Do(context.TODO(), req)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHTTPClientWithContextGetSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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

	response, err := client.Get(context.TODO(), server.URL, headers)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientWithContextPostSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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

	response, err := client.Post(context.TODO(), server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientWithContextDeleteSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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

	response, err := client.Delete(context.TODO(), server.URL, headers)
	require.NoError(t, err, "should not have failed to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientWithContextPutSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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

	response, err := client.Put(context.TODO(), server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientWithContextPatchSuccessWithTODOContext(t *testing.T) {
	client := NewHTTPClientWithContext(10 * time.Millisecond)

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

	response, err := client.Patch(context.TODO(), server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

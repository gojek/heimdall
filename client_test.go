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
	config := Config{
		timeoutInSeconds: 10,
	}

	client := NewHTTPClient(config)

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

func TestHTTPClientPostSuccess(t *testing.T) {
	config := Config{
		timeoutInSeconds: 10,
	}

	client := NewHTTPClient(config)

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

func TestHTTPClientDeleteSuccess(t *testing.T) {
	config := Config{
		timeoutInSeconds: 10,
	}

	client := NewHTTPClient(config)

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

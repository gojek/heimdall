package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

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

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHTTPClientGetSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

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

func TestHTTPClientPostSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
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

func TestHTTPClientDeleteSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

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

func TestHTTPClientPutSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
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

func TestHTTPClientPatchSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientGetRetriesOnFailure(t *testing.T) {
	count := 0
	noOfRetries := 3
	noOfCalls := noOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err, "should have failed to make GET request")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))

	assert.Equal(t, noOfCalls, count)
}

func BenchmarkHTTPClientGetRetriesOnFailure(b *testing.B) {
	noOfRetries := 3
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for b.Loop() {
		_, _ = client.Get(server.URL, http.Header{})
	}
}

func TestHTTPClientPostRetriesOnFailure(t *testing.T) {
	count := 0
	noOfRetries := 3
	noOfCalls := noOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Post(server.URL, strings.NewReader("a=1"), http.Header{})
	require.NoError(t, err, "should have failed to make GET request")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))

	assert.Equal(t, noOfCalls, count)
}

func BenchmarkHTTPClientPostRetriesOnFailure(b *testing.B) {
	noOfRetries := 3
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for b.Loop() {
		_, _ = client.Post(server.URL, strings.NewReader("a=1"), http.Header{})
	}
}

func TestHTTPClientGetReturnsNoErrorsIfRetriesFailWith5xx(t *testing.T) {
	count := 0
	noOfRetries := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)

	require.Equal(t, noOfRetries+1, count)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))
}

func TestHTTPClientGetReturnsNoErrorsIfRetrySucceeds(t *testing.T) {
	count := 0
	countWhenCallSucceeds := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

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

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err, "should not have failed to make GET request")

	require.Equal(t, countWhenCallSucceeds+1, count)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "{ \"response\": \"success\" }", respBody(t, response))
}

func TestHTTPClientGetReturnsErrorOnClientCallFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	server.URL = "" // Invalid URL to simulate client.Do error
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.Error(t, err, "should have failed to make GET request")

	require.Nil(t, response)

	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestHTTPClientGetReturnsNoErrorOn5xxFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)

}

func TestHTTPClientGetReturnsErrorOnFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	response, err := client.Get("url_doenst_exist", http.Header{})
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
	assert.Nil(t, response)
}

func TestPluginMethodsCalled(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))
	mockPlugin := &MockPlugin{}
	client.AddPlugin(mockPlugin)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	mockPlugin.On("OnRequestStart", mock.Anything)
	mockPlugin.On("OnRequestEnd", mock.Anything, mock.Anything)

	_, err := client.Get(server.URL, http.Header{})

	require.NoError(t, err)
	mockPlugin.AssertNumberOfCalls(t, "OnRequestStart", 1)
	pluginRequest, ok := mockPlugin.Calls[0].Arguments[0].(*http.Request)
	require.True(t, ok)
	assert.Equal(t, http.MethodGet, pluginRequest.Method)
	assert.Equal(t, server.URL, pluginRequest.URL.String())

	mockPlugin.AssertNumberOfCalls(t, "OnRequestEnd", 1)
	pluginResponse, ok := mockPlugin.Calls[1].Arguments[1].(*http.Response)
	require.True(t, ok)
	assert.Equal(t, http.StatusOK, pluginResponse.StatusCode)
}

func TestPluginErrorMethodCalled(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))
	mockPlugin := &MockPlugin{}
	client.AddPlugin(mockPlugin)

	mockPlugin.On("OnRequestStart", mock.Anything)
	mockPlugin.On("OnError", mock.Anything, mock.Anything)

	serverURL := "does_not_exist"
	client.Get(serverURL, http.Header{})

	mockPlugin.AssertNumberOfCalls(t, "OnRequestStart", 1)
	pluginRequest, ok := mockPlugin.Calls[0].Arguments[0].(*http.Request)
	require.True(t, ok)
	assert.Equal(t, http.MethodGet, pluginRequest.Method)
	assert.Equal(t, serverURL, pluginRequest.URL.String())

	mockPlugin.AssertNumberOfCalls(t, "OnError", 1)
	err, ok := mockPlugin.Calls[1].Arguments[1].(error)
	require.True(t, ok)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Set("foo", "bar")
	return c.client.Do(request)
}

func TestCustomHTTPClientHeaderSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithHTTPClient(&myHTTPClient{
			client: http.Client{Timeout: 25 * time.Millisecond}}),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "bar", r.Header.Get("foo"))
		assert.NotEqual(t, "baz", r.Header.Get("foo"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, err := client.Do(req)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHTTPClientContextTimeout(t *testing.T) {
	client := NewClient(WithHTTPTimeout(1000 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		time.Sleep(100 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxt, http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")
	response, err := client.Do(req)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, &url.Error{Op: "Get", URL: server.URL, Err: context.DeadlineExceeded}, err)
	require.Nil(t, response)
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer response.Body.Close()
	}

	respBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")

	return string(respBody)
}

func TestHTTPClientDoContextCancelledDuringRetry(t *testing.T) {
	noOfRetries := 3
	backoffInterval := 100 * time.Millisecond
	maximumJitterInterval := 10 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	count := atomic.Int32{}
	count.Store(0)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		count.Add(1)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	// Cancel the context after a short delay to simulate context cancellation during sleep
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err = client.Do(req)
	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Less(t, count.Load(), int32(noOfRetries+1), "should not have completed all retries due to context cancellation")
}

func TestHTTPClientDoContextCancelledBeforeRetry(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(3),
		WithRetrier(heimdall.NewRetrierFunc(func(retry int) time.Duration {
			assert.Fail(t, "should not have retrier func due to context cancellation")
			return 0
		})),
	)
	ctx, cancel := context.WithCancel(context.Background())

	count := atomic.Int32{}
	count.Store(0)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		cancel() // Cancel immediately
		count.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), context.Canceled.Error())
	assert.Equal(t, int32(1), count.Load())
}

func TestHTTPClientDoContextTimeoutDuringRetry(t *testing.T) {
	noOfRetries := 3
	backoffInterval := 100 * time.Millisecond
	maximumJitterInterval := 10 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	count := atomic.Int32{}
	count.Store(0)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		count.Add(1)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), context.DeadlineExceeded.Error())
	assert.Less(t, count.Load(), int32(noOfRetries+1), "should not have completed all retries due to context timeout")
}

func TestHTTPClientMultiRetryOnTimeout(t *testing.T) {
	noOfRetries := 3
	backoffInterval := 4 * time.Millisecond
	maximumJitterInterval := 2 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(5*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	mutliErr, ok := err.(interface{ Unwrap() []error })
	require.True(t, ok)
	errs := mutliErr.Unwrap()
	assert.Len(t, errs, noOfRetries+1)
	for _, e := range errs {
		var urlErr *url.Error
		require.True(t, errors.As(e, &urlErr))
		assert.Equal(t, "Get", urlErr.Op)
		assert.Equal(t, server.URL, urlErr.URL)
		assert.ErrorIs(t, urlErr.Err, context.DeadlineExceeded)
	}
}

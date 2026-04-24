package internal

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetRequestGetBodyReturnsNilForNilRequest(t *testing.T) {
	t.Parallel()

	require.Nil(t, SetRequestGetBody(nil))
}

func TestSetRequestGetBodyReturnsNilWhenBodyIsNil(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	require.Nil(t, err)

	require.Nil(t, SetRequestGetBody(req))
	assert.Nil(t, req.GetBody)
}

func TestSetRequestGetBodyKeepsExistingGetBody(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(strings.NewReader("payload")))
	require.Nil(t, err)

	called := 0
	req.GetBody = func() (io.ReadCloser, error) {
		called++
		return io.NopCloser(strings.NewReader("existing")), nil
	}

	require.Nil(t, SetRequestGetBody(req))
	assert.Zero(t, called)

	body, err := req.GetBody()
	require.Nil(t, err)

	assert.Equal(t, "existing", readBody(t, body))
	assert.Equal(t, 1, called)
	assert.Equal(t, "payload", readBody(t, req.Body))
}

func TestSetRequestGetBodySupportsNoBody(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", http.NoBody)
	require.Nil(t, err)

	require.Nil(t, SetRequestGetBody(req))
	require.NotNil(t, req.GetBody)

	body, err := req.GetBody()
	require.Nil(t, err)

	assert.Equal(t, http.NoBody, body)
	assert.Equal(t, "", readBody(t, body))
}

func TestSetRequestGetBodyMakesRequestBodyReplayable(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(strings.NewReader("payload")))
	require.Nil(t, err)

	require.Nil(t, SetRequestGetBody(req))
	require.NotNil(t, req.GetBody)

	assert.Equal(t, "payload", readBody(t, req.Body))

	firstReplay, err := req.GetBody()
	require.Nil(t, err)
	assert.Equal(t, "payload", readBody(t, firstReplay))

	secondReplay, err := req.GetBody()
	require.Nil(t, err)
	assert.Equal(t, "payload", readBody(t, secondReplay))
}

func TestSetRequestGetBodyReturnsErrorWhenBodyReadFails(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	require.Nil(t, err)
	req.Body = errReadCloser{err: errors.New("read failed")}

	require.EqualError(t, SetRequestGetBody(req), "read failed")
	assert.Nil(t, req.GetBody)
}

func TestCloneRequestReturnsNilWhenRequestIsNil(t *testing.T) {
	t.Parallel()

	cloned, err := CloneRequest(nil, func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("payload")), nil
	})

	require.Nil(t, err)
	assert.Nil(t, cloned)
}

func TestCloneRequestReturnsOriginalRequestWhenGetBodyIsNil(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(strings.NewReader("payload")))
	require.Nil(t, err)

	cloned, err := CloneRequest(req, nil)

	require.Nil(t, err)
	assert.Same(t, req, cloned)
	assert.Equal(t, "payload", readBody(t, cloned.Body))
}

func TestCloneRequestReturnsErrorWhenGetBodyFails(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(strings.NewReader("payload")))
	require.Nil(t, err)

	cloned, err := CloneRequest(req, func() (io.ReadCloser, error) {
		return nil, errors.New("clone body failed")
	})

	require.EqualError(t, err, "clone body failed")
	assert.Nil(t, cloned)
}

func TestCloneRequestReturnsClonedRequestWithFreshBodyAndPreservedContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), ctxKey("key"), "value")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", io.NopCloser(strings.NewReader("original-body")))
	require.Nil(t, err)
	req.Header.Set("X-Request-Id", "123")

	cloned, err := CloneRequest(req, func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("cloned-body")), nil
	})

	require.Nil(t, err)
	require.NotNil(t, cloned)

	assert.NotSame(t, req, cloned)
	assert.Equal(t, req.Method, cloned.Method)
	assert.Equal(t, req.URL.String(), cloned.URL.String())
	assert.Equal(t, "123", cloned.Header.Get("X-Request-Id"))
	assert.Equal(t, "value", cloned.Context().Value(ctxKey("key")))
	assert.Equal(t, "original-body", readBody(t, req.Body))
	assert.Equal(t, "cloned-body", readBody(t, cloned.Body))
}

func readBody(t *testing.T, body io.ReadCloser) string {
	t.Helper()

	defer body.Close()

	data, err := io.ReadAll(body)
	require.Nil(t, err)

	return string(data)
}

type ctxKey string
type errReadCloser struct {
	err error
}

func (e errReadCloser) Read(_ []byte) (int, error) {
	return 0, e.err
}

func (e errReadCloser) Close() error {
	return nil
}

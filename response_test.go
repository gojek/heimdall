package heimdall

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCodeOfResponse(t *testing.T) {
	response := Response{
		statusCode: http.StatusForbidden,
		body:       []byte(`hello`),
	}

	assert.Equal(t, http.StatusForbidden, response.StatusCode())
}

func TestBodyOfResponse(t *testing.T) {
	response := Response{
		statusCode: http.StatusForbidden,
		body:       []byte(`hello`),
	}

	assert.Equal(t, []byte(`hello`), response.Body())
}

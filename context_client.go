package heimdall

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gojektech/valkyrie"
	"github.com/pkg/errors"
)

type httpClientWithContext struct {
	client *http.Client

	retryCount int
	retrier    Retriable
}

func (c *httpClientWithContext) Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

func (c *httpClientWithContext) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var response *http.Response
	var contextCancelled bool = false

	multiErr := &valkyrie.MultiError{}

	for i := 0; i <= c.retryCount; i++ {
		var err error
		response, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				contextCancelled = true
			}

			multiErr.Push(err.Error())
			if contextCancelled {
				break
			}
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		if response.StatusCode >= http.StatusInternalServerError {
			multiErr.Push(fmt.Sprintf("server error: %d", response.StatusCode))

			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			fmt.Println("R: ", response.StatusCode)
			continue
		}

		multiErr = &valkyrie.MultiError{} // Clear errors if any iteration succeeds
		break
	}

	return response, multiErr.HasError()
}

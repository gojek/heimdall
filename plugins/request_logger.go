package plugins

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gojektech/heimdall/v6"
)

type ctxKey string

const reqTime ctxKey = "request_time_start"

type requestLogger struct {
	out    io.Writer
	errOut io.Writer
}

// NewRequestLogger returns a new instance of a Heimdall request logger plugin
// out and errOut are the streams where standard and error logs are written respectively
// If given as nil, `out` takes the default value of `os.StdOut`
// and errOut takes the default value of `os.StdErr`
func NewRequestLogger(out io.Writer, errOut io.Writer) heimdall.Plugin {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stderr
	}
	return &requestLogger{
		out:    out,
		errOut: errOut,
	}
}

func (rl *requestLogger) OnRequestStart(req *http.Request) {
	ctx := context.WithValue(req.Context(), reqTime, time.Now())
	*req = *(req.WithContext(ctx))
}

func (rl *requestLogger) OnRequestEnd(req *http.Request, res *http.Response) {
	reqDurationMs := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	statusCode := res.StatusCode
	fmt.Fprintf(rl.out, "%s %s %s %d [%dms]\n", time.Now().Format("02/Jan/2006 03:04:05"), method, url, statusCode, reqDurationMs)
}

func (rl *requestLogger) OnError(req *http.Request, err error) {
	reqDurationMs := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	fmt.Fprintf(rl.errOut, "%s %s %s [%dms] ERROR: %v\n", time.Now().Format("02/Jan/2006 03:04:05"), method, url, reqDurationMs, err)
}

func getRequestDuration(ctx context.Context) time.Duration {
	now := time.Now()
	start := ctx.Value(reqTime)
	if start == nil {
		return 0
	}
	startTime, ok := start.(time.Time)
	if !ok {
		return 0
	}
	return now.Sub(startTime)
}

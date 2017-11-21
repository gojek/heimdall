## Heimdall
![Heimdall](https://i.stack.imgur.com/3eLbk.png)

[![Build Status](https://travis-ci.org/gojek-engineering/heimdall.svg?branch=master)](https://travis-ci.org/gojek-engineering/heimdall)

### Yet another Golang HTTP Client

- Provides wrapper over HTTPClient with timeouts and a fluent interface for calling HTTP Methods
- Provides Hystrix HTTP client with fluent interface
- Supports Synchronous in-memory retries

How to use this library:

- Add this library as a dependency to your glide.yaml file, and preferably fix a version

# heimdall
--
    import "github.com/gojek-engineering/heimdall"


## Usage

#### type Backoff

```go
type Backoff interface {
	Next(retry int) time.Duration
}
```

Backoff interface defines contract for backoff strategies

#### func  NewConstantBackoff

```go
func NewConstantBackoff(backoffInterval int64) Backoff
```
NewConstantBackoff returns an instance of ConstantBackoff

#### func  NewExponentialBackoff

```go
func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64) Backoff
```
NewExponentialBackoff returns an instance of ExponentialBackoff

#### type Client

```go
type Client interface {
	Get(url string, headers http.Header) (Response, error)
	Post(url string, body io.Reader, headers http.Header) (Response, error)
	Put(url string, body io.Reader, headers http.Header) (Response, error)
	Patch(url string, body io.Reader, headers http.Header) (Response, error)
	Delete(url string, headers http.Header) (Response, error)

	SetRetryCount(count int)
	SetRetrier(retrier Retriable)
}
```

Client Is a generic HTTP client interface

#### func  NewHTTPClient

```go
func NewHTTPClient(timeoutInMilliseconds int) Client
```
NewHTTPClient returns a new instance of HTTPClient

#### func  NewHystrixHTTPClient

```go
func NewHystrixHTTPClient(httpClient *http.Client, hystrixConfig HystrixConfig) Client
```
NewHystrixHTTPClient returns a new instance of HystrixHTTPClient

#### type HystrixCommandConfig

```go
type HystrixCommandConfig struct {
	Timeout                int
	MaxConcurrentRequests  int
	RequestVolumeThreshold int
	SleepWindow            int
	ErrorPercentThreshold  int
}
```

HystrixCommandConfig takes the hystrix config values

#### type HystrixConfig

```go
type HystrixConfig struct {
}
```

HystrixConfig is used to pass configurations for Hystrix

#### func  NewHystrixConfig

```go
func NewHystrixConfig(commandName string, commandConfig HystrixCommandConfig) HystrixConfig
```
NewHystrixConfig should be used to give hystrix commandName and config

#### type Response

```go
type Response struct {
}
```

Response encapsulates details of a http response

#### func (Response) Body

```go
func (hr Response) Body() []byte
```
Body returns body in bytes of a http request

#### func (Response) StatusCode

```go
func (hr Response) StatusCode() int
```
StatusCode returns status code of a http request

#### type Retriable

```go
type Retriable interface {
	NextInterval(retry int) time.Duration
}
```

Retriable defines contract for retriers to implement

#### func  NewNoRetrier

```go
func NewNoRetrier() Retriable
```
NewNoRetrier returns a null object for retriable

#### func  NewRetrier

```go
func NewRetrier(backoff Backoff) Retriable
```
NewRetrier returns retrier with some backoff strategy

TODO:

- [ ] Support Connection Pooling at transport layer
- [ ] Fallback support for hystrix client
- [ ] Instrumentation of these calls using HTTPClient


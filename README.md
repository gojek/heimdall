# Heimdall

<p align="center"><img src="doc/logo.png" width="360"></p>
<p align="center">
  <a href="https://travis-ci.org/gojektech/heimdall"><img src="https://travis-ci.org/gojektech/heimdall.svg?branch=master" alt="Build Status"></img></a>
  <a href="https://goreportcard.com/report/github.com/gojektech/heimdall"><img src="https://goreportcard.com/badge/github.com/gojektech/heimdall"></img></a>
  <a href="https://golangci.com"><img src="https://golangci.com/badges/github.com/gojektech/heimdall.svg"></img></a>
</p>

## Description

Heimdall is an HTTP client that helps your application make a large number of requests, at scale. With Heimdall, you can:
- Use a [hystrix-like](https://github.com/afex/hystrix-go) circuit breaker to control failing requests
- Add synchronous in-memory retries to each request, with the option of setting your own retrier strategy
- Create clients with different timeouts for every request

All HTTP methods are exposed as a fluent interface.

## Installation
```
go get -u github.com/gojektech/heimdall
```

## Usage

### Making a simple `GET` request
The below example will print the contents of the google home page:

```go
// Create a new HTTP client with a default timeout
timeout := 1000 * time.Millisecond
client := heimdall.NewHTTPClient(timeout)

// Use the clients GET method to create and execute the request
res, err := client.Get("http://google.com", nil)
if err != nil{
	panic(err)
}

// Heimdall returns the standard *http.Response object
body, err := ioutil.ReadAll(res.Body)
fmt.Println(string(body))
```

You can also use the `*http.Request` object with the `http.Do` interface :

```go
timeout := 1000 * time.Millisecond
client := heimdall.NewHTTPClient(timeout)

// Create an http.Request instance
req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
// Call the `Do` method, which has a similar interface to the `http.Do` method
res, err := client.Do(req)
if err != nil {
	panic(err)
}

body, err := ioutil.ReadAll(res.Body)
fmt.Println(string(body))
```

### Creating a hystrix-like circuit breaker

You can use the `NewHystrixHTTPClient` function to create a client wrapped in a hystrix-like circuit breaker:

```go
// Create a new hystrix config, and input the command name, along with other required options
hystrixConfig := heimdall.NewHystrixConfig("google_get_request", heimdall.HystrixCommandConfig{
    ErrorPercentThreshold : 20,
    MaxConcurrentRequests: 30,
    Timeout: 1000,
})
// Create a new hystrix-wrapped HTTP client
client := heimdall.NewHystrixHTTPClient(1000, hystrixConfig)

// The rest is the same as the previous example
```

In the above example, there are two timeout values used: one for the hystrix configuration, and one for the HTTP client configuration. The former determines the time at which hystrix should register an error, while the latter determines when the client itself should return a timeout error. Unless you have any special requirements, both of these would have the same values.

### Creating a hystrix-like circuit breaker with fallbacks

You can use the `NewHystrixHTTPClient` function to create a client wrapped in a hystrix-like circuit breaker by passing in your own custom fallbacks:

The fallback function will trigger when your code returns an error, or whenever it is unable to complete based on a variety of [health checks](https://github.com/Netflix/Hystrix/wiki/How-it-Works).

**How your fallback function should look like**
you should pass in a function whose signature looks like following
```go
func(err error) error {
    // your logic for handling the error/outage condition
    return err
}
```


**Example**
```go
// Create a new hystrix config, and input the command name, along with other required options
hystrixConfig := heimdall.NewHystrixConfig("post_to_channel_one", heimdall.HystrixCommandConfig{
    ErrorPercentThreshold : 20,
    MaxConcurrentRequests: 30,
    Timeout: 1000,
})
// Create a new fallback function
fallbackFn := func(err error) error {
    _, err := http.Post("post_to_channel_two")
    return err
}

hystrixConfig := heimdall.NewHystrixConfig("MyCommand", heimdall.HystrixCommandConfig{
	Timeout:                1100,
	MaxConcurrentRequests:  100,
	ErrorPercentThreshold:  25,
	SleepWindow:            10,
	RequestVolumeThreshold: 10,
	fallbackFunc: fallbackFn,
})

// Create a new hystrix-wrapped HTTP client with the fallbackFunc as fall-back function
client := heimdall.NewHystrixHTTPClient(1000, hystrixConfig)

// The rest is the same as the previous example
```

In the above example, the `fallbackFunc` is a function which posts to channel two in case posting to channel one fails.

### Creating an HTTP client with a retry mechanism

```go
// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
backoff := heimdall.NewConstantBackoff(500)
// Create a new retry mechanism with the backoff
retrier := heimdall.NewRetrier(backoff)

client := heimdall.NewHTTPClient(1000)
// Set the retry mechanism for the client, and the number of times you would like to retry
client.SetRetrier(retrier)
client.SetRetryCount(4)

// The rest is the same as the first example
```

This will create an HTTP client which will retry every `500` milliseconds incase the request fails. THe library also comes with an [Exponential Backoff](https://www.godoc.org/github.com/gojektech/heimdall#NewExponentialBackoff)

### Custom retry mechanisms

Heimdall supports custom retry strategies. To do this, you will have to implement the `Backoff` interface:

```go
type Backoff interface {
	Next(retry int) time.Duration
}
```

Let's see an example of creating a client with a linearly increasing backoff time:

First, create the backoff mechanism:

```go
type linearBackoff struct {
	backoffInterval int
}

func (lb *linearBackoff) Next(retry int) time.Duration{
	if retry <= 0 {
		return 0 * time.Millisecond
	}
	return time.Duration(retry * lb.backoffInterval) * time.Millisecond
}
```

This will create a backoff mechanism, where the retry time will increase linearly for each retry attempt. We can use this to create the client, just like the last example:

```go
backoff := &linearBackoff{100}
retrier := heimdall.NewRetrier(backoff)

client := heimdall.NewHTTPClient(1000)
client.SetRetrier(retrier)
client.SetRetryCount(4)

// The rest is the same as the first example
```

## Documentation

Further documentation can be found on [godoc.org](https://www.godoc.org/github.com/gojektech/heimdall)

## License

```
Copyright 2018, GO-JEK Tech (http://gojek.tech)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

# Heimdall
[![Build Status](https://travis-ci.org/gojektech/heimdall.svg?branch=master)](https://travis-ci.org/gojektech/heimdall)

![Heimdall Logo](doc/logo.png)

## Description

Heimdall is an HTTP client that helps your application make a large number of requests, at scale. With Heimdall, you can:
- Use a [hystrix-like](https://github.com/afex/hystrix-go) circuit breaker to control failing requests
- Add synchronous in-memory retries to each request, with the option of setting your own retrier strategy
- Create clients with different timeouts for every request

All HTTP methods are exposed as a fluent interface. The 

## Installation 
```
go get -u github.com/gojektech/heimdall
```

## Usage

### Making a simple `GET` request
The below example will print the contents of the google home page:

```go
// Create a new HTTP client with a default timeout
client := heimdall.NewHTTPClient(1000)

// Use the clients GET method to create and execute the request
res, err := client.Get("http://google.com", nil)
if err != nil{
	panic(err)
}

// The heimdall response object comes with handy methods to obtain the contents of the reponse
// In this, case we can directly get the bytes of the response body using the `Body` method
fmt.Println(string(res.Body()))
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

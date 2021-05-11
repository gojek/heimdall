# Heimdall

<p align="center"><img src="doc/heimdall-logo.png" width="360"></p>
<p align="center">
  <a href="https://travis-ci.com/gojek/heimdall"><img src="https://travis-ci.com/gojek/heimdall.svg?branch=master" alt="Build Status"></img></a>
  <a href="https://goreportcard.com/report/github.com/gojek/heimdall"><img src="https://goreportcard.com/badge/github.com/gojek/heimdall"></img></a>
  <a href="https://golangci.com"><img src="https://golangci.com/badges/github.com/gojek/heimdall.svg"></img></a>
  <a href="https://coveralls.io/github/gojek/heimdall?branch=master"><img src="https://coveralls.io/repos/github/gojek/heimdall/badge.svg?branch=master"></img></a>
</p>

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
  + [Making a simple `GET` request](#making-a-simple-get-request)
  + [Creating a hystrix-like circuit breaker](#creating-a-hystrix-like-circuit-breaker)
  + [Creating a hystrix-like circuit breaker with fallbacks](#creating-a-hystrix-like-circuit-breaker-with-fallbacks)
  + [Creating an HTTP client with a retry mechanism](#creating-an-http-client-with-a-retry-mechanism)
  + [Custom retry mechanisms](#custom-retry-mechanisms)
  + [Custom HTTP clients](#custom-http-clients)
* [Plugins](#plugins)
* [Documentation](#documentation)
* [FAQ](#faq)
* [License](#license)

## Description

Heimdall is an HTTP client that helps your application make a large number of requests, at scale. With Heimdall, you can:
- Use a [hystrix-like](https://github.com/afex/hystrix-go) circuit breaker to control failing requests
- Add synchronous in-memory retries to each request, with the option of setting your own retrier strategy
- Create clients with different timeouts for every request

All HTTP methods are exposed as a fluent interface.

## Installation
```
go get -u github.com/gojek/heimdall/v7
```

## Usage

### Importing the package

This package can be used by adding the following import statement to your `.go` files.

```go
import "github.com/gojek/heimdall/v7/httpclient" 
```

### Making a simple `GET` request
The below example will print the contents of the google home page:

```go
// Create a new HTTP client with a default timeout
timeout := 1000 * time.Millisecond
client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

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
client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

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

To import hystrix package of heimdall.

```go
import "github.com/gojek/heimdall/v7/hystrix"
```

You can use the `hystrix.NewClient` function to create a client wrapped in a hystrix-like circuit breaker:

```go
// Create a new hystrix-wrapped HTTP client with the command name, along with other required options
client := hystrix.NewClient(
	hystrix.WithHTTPTimeout(10 * time.Millisecond),
	hystrix.WithCommandName("google_get_request"),
	hystrix.WithHystrixTimeout(1000 * time.Millisecond),
	hystrix.WithMaxConcurrentRequests(30),
	hystrix.WithErrorPercentThreshold(20),
	hystrix.WithStatsDCollector("localhost:8125", "myapp.hystrix"),
)

// The rest is the same as the previous example
```

In the above example, there are two timeout values used: one for the hystrix configuration, and one for the HTTP client configuration. The former determines the time at which hystrix should register an error, while the latter determines when the client itself should return a timeout error. Unless you have any special requirements, both of these would have the same values.

You can choose to export hystrix metrics to a statsD collector with the `hystrix.WithStatsDCollector(<statsd addr>, <metrics-prefix>)` option when initializing the client as shown above.

### Creating a hystrix-like circuit breaker with fallbacks

You can use the `hystrix.NewClient` function to create a client wrapped in a hystrix-like circuit breaker by passing in your own custom fallbacks:

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
// Create a new fallback function
fallbackFn := func(err error) error {
    _, err := http.Post("post_to_channel_two")
    return err
}

timeout := 10 * time.Millisecond

// Create a new hystrix-wrapped HTTP client with the fallbackFunc as fall-back function
client := hystrix.NewClient(
	hystrix.WithHTTPTimeout(timeout),
	hystrix.WithCommandName("MyCommand"),
	hystrix.WithHystrixTimeout(1100 * time.Millisecond),
	hystrix.WithMaxConcurrentRequests(100),
	hystrix.WithErrorPercentThreshold(20),
	hystrix.WithSleepWindow(10),
	hystrix.WithRequestVolumeThreshold(10),
	hystrix.WithFallbackFunc(fallbackFn),
})

// The rest is the same as the previous example
```

In the above example, the `fallbackFunc` is a function which posts to channel two in case posting to channel one fails.

### Creating an HTTP client with a retry mechanism

```go
// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
backoffInterval := 2 * time.Millisecond
// Define a maximum jitter interval. It must be more than 1*time.Millisecond
maximumJitterInterval := 5 * time.Millisecond

backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

// Create a new retry mechanism with the backoff
retrier := heimdall.NewRetrier(backoff)

timeout := 1000 * time.Millisecond
// Create a new client, sets the retry mechanism, and the number of times you would like to retry
client := httpclient.NewClient(
	httpclient.WithHTTPTimeout(timeout),
	httpclient.WithRetrier(retrier),
	httpclient.WithRetryCount(4),
)

// The rest is the same as the first example
```
Or create client with exponential backoff

```go
// First set a backoff mechanism. Exponential Backoff increases the backoff at a exponential rate

initalTimeout := 2*time.Millisecond            // Inital timeout
maxTimeout := 9*time.Millisecond               // Max time out
exponentFactor := 2                            // Multiplier
maximumJitterInterval := 2*time.Millisecond    // Max jitter interval. It must be more than 1*time.Millisecond

backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, exponentFactor, maximumJitterInterval)

// Create a new retry mechanism with the backoff
retrier := heimdall.NewRetrier(backoff)

timeout := 1000 * time.Millisecond
// Create a new client, sets the retry mechanism, and the number of times you would like to retry
client := httpclient.NewClient(
	httpclient.WithHTTPTimeout(timeout),
	httpclient.WithRetrier(retrier),
	httpclient.WithRetryCount(4),
)

// The rest is the same as the first example
```

This will create an HTTP client which will retry every `500` milliseconds incase the request fails. The library also comes with an [Exponential Backoff](https://pkg.go.dev/github.com/gojek/heimdall#NewExponentialBackoff)

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

timeout := 1000 * time.Millisecond
// Create a new client, sets the retry mechanism, and the number of times you would like to retry
client := httpclient.NewClient(
	httpclient.WithHTTPTimeout(timeout),
	httpclient.WithRetrier(retrier),
	httpclient.WithRetryCount(4),
)

// The rest is the same as the first example
```

Heimdall also allows you to simply pass a function that returns the retry timeout. This can be used to create the client, like:
```go
linearRetrier := NewRetrierFunc(func(retry int) time.Duration {
	if retry <= 0 {
		return 0 * time.Millisecond
	}
	return time.Duration(retry) * time.Millisecond
})

timeout := 1000 * time.Millisecond
client := httpclient.NewClient(
	httpclient.WithHTTPTimeout(timeout),
	httpclient.WithRetrier(linearRetrier),
	httpclient.WithRetryCount(4),
)
```

### Custom HTTP clients

Heimdall supports custom HTTP clients. This is useful if you are using a client imported from another library and/or wish to implement custom logging, cookies, headers etc for each request that you make with your client.

Under the hood, the `httpClient` struct now accepts `Doer`, which is the standard interface implemented by HTTP clients (including the standard library's `net/*http.Client`)

Let's say we wish to add authorization headers to all our requests.

We can define our client `myHTTPClient`

```go
type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth("username", "passwd")
	return c.client.Do(request)
}
```

And set this with `httpclient.NewClient(httpclient.WithHTTPClient(&myHTTPClient{client: http.DefaultClient}))`

Now, each sent request will have the `Authorization` header to use HTTP basic authentication with the provided username and password.

This can be done for the hystrix client as well

```go
client := httpclient.NewClient(
	httpclient.WithHTTPClient(&myHTTPClient{
		client: http.Client{Timeout: 25 * time.Millisecond},
	}),
)

// The rest is the same as the first example
```

## Plugins

To add a plugin to an existing client, use the `AddPlugin` method of the client. 

An example, with the [request logger plugin](/plugins/request_logger.go):

```go
// import "github.com/gojek/heimdall/v7/plugins"

client := heimdall.NewHTTPClient(timeout)
requestLogger := plugins.NewRequestLogger(nil, nil)
client.AddPlugin(requestLogger)
// use the client as before

req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
res, err := client.Do(req)
if err != nil {
	panic(err)
}
// This will log:
//23/Jun/2018 12:48:04 GET http://google.com 200 [412ms]
// to STDOUT
```

A plugin is an interface whose methods get called during key events in a requests lifecycle:

- `OnRequestStart` is called just before the request is made
- `OnRequestEnd` is called once the request has successfully executed
- `OnError` is called is the request failed

Each method is called with the request object as an argument, with `OnRequestEnd`, and `OnError` additionally being called with the response and error instances respectively.
For a simple example on how to write plugins, look at the [request logger plugin](/plugins/request_logger.go).

## Documentation

Further documentation can be found on [pkg.go.dev](https://pkg.go.dev/github.com/gojek/heimdall/v7)

## FAQ

**Can I replace the standard Go HTTP client with Heimdall?**

Yes, you can. Heimdall implements the standard [HTTP Do](https://golang.org/pkg/net/http/#Client.Do) method, along with [useful wrapper methods](https://golang.org/pkg/net/http/#Client.Do) that provide all the functionality that a regular Go HTTP client provides.

---

**When should I use Heimdall?**

If you are making a large number of HTTP requests, or if you make requests among multiple distributed nodes, and wish to make your systems more fault tolerant, then Heimdall was made for you.

Heimdall makes use of [multiple mechanisms](https://medium.com/@sohamkamani/how-to-handle-microservice-communication-at-scale-a6fb0ee0ed7) to make HTTP requests more fault tolerant:
1. Retries - If a request fails, Heimdall retries behind the scenes, and returns the result if one of the retries are successful.
2. Circuit breaking - If Heimdall detects that too many of your requests are failing, or that the number of requests sent are above a configured threshold, then it "opens the circuit" for a short period of time, which prevents any more requests from being made. _This gives your downstream systems time to recover._

---

**So does this mean that I shouldn't use Heimdall for small scale applications?**

Although Heimdall was made keeping large scale systems in mind, it's interface is simple enough to be used for any type of systems. In fact, we use it for our pet projects as well. Even if you don't require retries or circuit breaking features, the [simpler HTTP client](https://github.com/gojek/heimdall#making-a-simple-get-request) provides sensible defaults with a simpler interface, and can be upgraded easily should the need arise.

---

**Can I contribute to make Heimdall better?**

[Please do!](https://github.com/gojek/heimdall/blob/master/CONTRIBUTING.md) We are looking for any kind of contribution to improve Heimdalls core funtionality and documentation. When in doubt, make a PR!

## License

```
Copyright 2018-2020, GO-JEK Tech (http://gojek.tech)

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

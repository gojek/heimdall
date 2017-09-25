## Heimdall
![Heimdall](https://i.stack.imgur.com/3eLbk.png)
### Yet another HTTP Client 

#### Supports Synchronous in-memory retries

How to use this library:

- Add this library as a dependency to your glide.yaml file, and preferably fix a version.
- Use
```
func NewHystrixHTTPClient(httpClient *http.Client, hystrixConfig *HystrixConfig) Client
```

to create a new http client with hystrix.

- HystrixConfig consists of
```
type HystrixConfig struct {
	commandName   string
	commandConfig *hystrix.CommandConfig {
         Timeout                int
         MaxConcurrentRequests  int
         RequestVolumeThreshold int
         SleepWindow            int
         ErrorPercentThreshold  int
    }
}

```
- This http Client provides methods like GET, POST, PUT etc which can be called to make http requests.
```
(hhc *hystrixHTTPClient) Post(url string, body io.Reader, headers http.Header) (Response, error)
```

Things to add:

- [ ] Instrumentation of external calls
- [ ] Integration with go-worker for asynchronous retries


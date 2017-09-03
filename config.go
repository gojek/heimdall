package heimdall

// Config used to define the parameters for a client
type Config struct {
	timeoutInSeconds int
}

// NewConfig returns an instance of config with provided params
func NewConfig(timeoutInSeconds int) Config {
	return Config{
		timeoutInSeconds: timeoutInSeconds,
	}
}

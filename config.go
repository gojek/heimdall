package heimdall

type Config struct {
	timeoutInSeconds int
}

func NewConfig(timeoutInSeconds int) Config {
	return Config{
		timeoutInSeconds: timeoutInSeconds,
	}
}

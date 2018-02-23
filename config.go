package heimdall

import (
	"github.com/afex/hystrix-go/hystrix"
)

type fallbackFunction func(error) error

// HystrixConfig is used to pass configurations for Hystrix
type HystrixConfig struct {
	commandName   string
	commandConfig hystrix.CommandConfig
	fallbackFunc  fallbackFunction
}

// HystrixCommandConfig takes the hystrix config values
type HystrixCommandConfig struct {
	Timeout                int
	MaxConcurrentRequests  int
	RequestVolumeThreshold int
	SleepWindow            int
	ErrorPercentThreshold  int
}

// NewHystrixConfig should be used to give hystrix commandName and config
func NewHystrixConfig(commandName string, commandConfig HystrixCommandConfig, f ...fallbackFunction) HystrixConfig {
	var hystrixFallbackFunction fallbackFunction
	if len(f) > 0 {
		hystrixFallbackFunction = f[0]
	}
	return HystrixConfig{
		commandName: commandName,
		commandConfig: hystrix.CommandConfig{
			Timeout:                commandConfig.Timeout,
			MaxConcurrentRequests:  commandConfig.MaxConcurrentRequests,
			RequestVolumeThreshold: commandConfig.RequestVolumeThreshold,
			SleepWindow:            commandConfig.SleepWindow,
			ErrorPercentThreshold:  commandConfig.ErrorPercentThreshold,
		},
		fallbackFunc: hystrixFallbackFunction,
	}
}

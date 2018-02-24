package heimdall

import (
	"testing"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestNewHystrixConfig(t *testing.T) {

	hystrixCommandConfig := HystrixCommandConfig{
		Timeout:                1,
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 200,
		SleepWindow:            50,
		ErrorPercentThreshold:  20,
	}

	hystrixConfig := NewHystrixConfig("hello", hystrixCommandConfig)

	expectedHystrixCommandConfig := hystrix.CommandConfig{
		Timeout:                1,
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 200,
		SleepWindow:            50,
		ErrorPercentThreshold:  20,
	}

	assert.Equal(t, "hello", hystrixConfig.commandName)
	assert.Equal(t, expectedHystrixCommandConfig, hystrixConfig.commandConfig)
}

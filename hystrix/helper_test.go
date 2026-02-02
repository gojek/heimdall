package hystrix

import (
	"sync"

	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
)

type simpleMetricCollector struct {
	name  string
	mx    sync.RWMutex
	total metricCollector.MetricResult
}

func newSimpleMetricCollector(name string) *simpleMetricCollector {
	return &simpleMetricCollector{name: name}
}

func (c *simpleMetricCollector) Reset() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.total = metricCollector.MetricResult{}
}

func (c *simpleMetricCollector) Update(r metricCollector.MetricResult) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.total.Attempts += r.Attempts
	c.total.Errors += r.Errors
	c.total.Successes += r.Successes
	c.total.Failures += r.Failures
	c.total.Rejects += r.Rejects
	c.total.ShortCircuits += r.ShortCircuits
	c.total.Timeouts += r.Timeouts
	c.total.FallbackSuccesses += r.FallbackSuccesses
	c.total.FallbackFailures += r.FallbackFailures
	c.total.ContextCanceled += r.ContextCanceled
	c.total.ContextDeadlineExceeded += r.ContextDeadlineExceeded
	c.total.TotalDuration += r.TotalDuration
	c.total.RunDuration += r.RunDuration
	c.total.ConcurrencyInUse += r.ConcurrencyInUse
}

func (c *simpleMetricCollector) GetMetrics() metricCollector.MetricResult {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.total
}

type simpleMetricRegistry struct {
	collectors map[string]*simpleMetricCollector
	mx         sync.RWMutex
}

func newSimpleMetricRegistry() *simpleMetricRegistry {
	r := &simpleMetricRegistry{
		collectors: make(map[string]*simpleMetricCollector),
	}
	metricCollector.Registry.Register(r.Register)
	return r
}

func (r *simpleMetricRegistry) GetMetrics(name string) metricCollector.MetricResult {
	r.mx.RLock()
	defer r.mx.RUnlock()
	collector := r.collectors[name]
	if collector == nil {
		return metricCollector.MetricResult{}
	}

	return collector.GetMetrics()
}

func (r *simpleMetricRegistry) Register(name string) metricCollector.MetricCollector {
	r.mx.Lock()
	defer r.mx.Unlock()

	collector := newSimpleMetricCollector(name)
	r.collectors[name] = collector
	return collector
}

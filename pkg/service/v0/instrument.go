package service

import (
	"github.com/owncloud/ocis-hello/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next Greeter, metrics *metrics.Metrics) Greeter {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    Greeter
	metrics *metrics.Metrics
}

// Greet implements the Greeter interface.
func (i instrument) Greet(accountID, name string) string {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000

		i.metrics.Latency.WithLabelValues().Observe(us)
		i.metrics.Duration.WithLabelValues().Observe(v)
	}))

	defer timer.ObserveDuration()

	i.metrics.Counter.WithLabelValues().Inc()

	return i.next.Greet(accountID, name)
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Namespace defines the namespace for the defines metrics.
	Namespace = "ocis"

	// Subsystem defines the subsystem for the defines metrics.
	Subsystem = "hello"
)

// Metrics defines the available metrics of this service.
type Metrics struct {
	Counter  *prometheus.CounterVec
	Latency  *prometheus.SummaryVec
	Duration *prometheus.HistogramVec
}

// New initializes the available metrics.
func New() *Metrics {
	m := &Metrics{
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "greet_total",
			Help:      "How many greeting requests processed",
		}, []string{}),
		Latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "greet_latency_microseconds",
			Help:      "Greet request latencies in microseconds",
		}, []string{}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      "greet_duration_seconds",
			Help:      "Greet method request time in seconds",
		}, []string{}),
	}

	prometheus.Register(
		m.Counter,
	)

	prometheus.Register(
		m.Latency,
	)

	prometheus.Register(
		m.Duration,
	)

	return m
}

package svc

import (
	"context"

	"github.com/owncloud/ocis-hello/pkg/metrics"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next v0proto.HelloHandler, metrics *metrics.Metrics) v0proto.HelloHandler {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    v0proto.HelloHandler
	metrics *metrics.Metrics
}

// Greet implements the HelloHandler interface.
func (i instrument) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000

		i.metrics.Latency.WithLabelValues().Observe(us)
		i.metrics.Duration.WithLabelValues().Observe(v)
	}))

	defer timer.ObserveDuration()

	err := i.next.Greet(ctx, req, rsp)

	if err == nil {
		i.metrics.Counter.WithLabelValues().Inc()
	}

	return err
}

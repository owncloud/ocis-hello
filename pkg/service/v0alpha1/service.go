package v0alpha1svc

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type options struct {
	logger  log.Logger
	counter metrics.Counter
}

// Service describes the provided service.
type Service interface {
	Greet(ctx context.Context, name string) (string, error)
}

// New returns a new service with all of the expected middlewares wired in.
func New(opts ...Option) Service {
	o := new(options)

	for _, opt := range opts {
		opt(o)
	}

	svc := Basic()

	{
		svc = Logging(log.With(o.logger, "type", "service"))(svc)
		svc = Instrument(o.counter)(svc)
	}

	return svc
}

type Middleware func(Service) Service

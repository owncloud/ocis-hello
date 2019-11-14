package v0alpha1endpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opencensus"
	"github.com/owncloud/ocis-hello/pkg/service/v0alpha1"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

var (
	_ endpoint.Failer = GreetResponse{}
)

type options struct {
	service  v0alpha1svc.Service
	logger   log.Logger
	duration metrics.Histogram
}

type Set struct {
	GreetEndpoint endpoint.Endpoint
}

func (s Set) Greet(ctx context.Context, name string) (string, error) {
	resp, err := s.GreetEndpoint(ctx, GreetRequest{
		Name: name,
	})

	if err != nil {
		return "", err
	}

	response := resp.(GreetResponse)
	return response.Message, response.Err
}

func New(opts ...Option) Set {
	o := new(options)

	for _, opt := range opts {
		opt(o)
	}

	greet := MakeGreetEndpoint(o.service)

	{
		greet = ratelimit.NewErroringLimiter(
			rate.NewLimiter(
				rate.Every(time.Second),
				1,
			),
		)(greet)

		greet = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(
				gobreaker.Settings{},
			),
		)(greet)

		greet = opencensus.TraceEndpoint(
			"hello/greet",
		)(greet)

		greet = Logging(log.With(o.logger, "type", "endpoint"))(greet)
		greet = Instrument(o.duration.With("method", "greet"))(greet)
	}

	return Set{
		GreetEndpoint: greet,
	}
}

// MakeGreetEndpoint constructs a Greet endpoint wrapping the service.
func MakeGreetEndpoint(svc v0alpha1svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GreetRequest)
		res, err := svc.Greet(ctx, req.Name)

		return GreetResponse{
			Message: res,
			Err:     err,
		}, nil
	}
}

type GreetRequest struct {
	Name string
}

type GreetResponse struct {
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (r GreetResponse) Failed() error {
	return r.Err
}

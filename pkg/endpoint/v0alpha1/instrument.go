package v0alpha1endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

func Instrument(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			start := time.Now()
			res, err := next(ctx, request)

			duration.With(
				"success",
				fmt.Sprint(err == nil),
			).Observe(
				time.Since(start).Seconds(),
			)

			return res, err
		}
	}
}

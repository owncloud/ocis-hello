package v0alpha1endpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func Logging(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			start := time.Now()
			res, err := next(ctx, request)

			logger = log.With(logger,
				"endpoint", "Greet",
				"duration", time.Since(start),
			)

			if err != nil {
				level.Warn(logger).Log(
					"err", err,
				)
			} else {
				level.Debug(logger).Log()
			}

			return res, err
		}
	}
}

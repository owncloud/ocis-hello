package v0alpha1svc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Logging takes a logger as a dependency and returns a service middleware.
func Logging(logger log.Logger) Middleware {
	return func(next Service) Service {
		return logging{
			next:   next,
			logger: logger,
		}
	}
}

type logging struct {
	next   Service
	logger log.Logger
}

func (l logging) Greet(ctx context.Context, name string) (string, error) {
	start := time.Now()
	res, err := l.next.Greet(ctx, name)

	logger := log.With(l.logger,
		"method", "Greet",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "Failed to execute",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return res, err
}

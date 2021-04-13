package service

import (
	"context"
	"time"

	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// NewLogging returns a service that logs messages.
func NewLogging(next v0proto.HelloHandler, logger log.Logger) v0proto.HelloHandler {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   v0proto.HelloHandler
	logger log.Logger
}

// Greet implements the HelloHandler interface.
func (l logging) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	start := time.Now()
	err := l.next.Greet(ctx, req, rsp)

	logger := l.logger.With().
		Str("method", "Hello.Greet").
		Dur("duration", time.Since(start)).
		Logger()

	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to execute")
	} else {
		logger.Debug().
			Msg("")
	}

	return err
}

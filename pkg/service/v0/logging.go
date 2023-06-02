package service

import (
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// NewLogging returns a service that logs messages.
func NewLogging(next Greeter, logger log.Logger) Greeter {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   Greeter
	logger log.Logger
}

// Greet implements the Greeter interface.
func (l logging) Greet(accountID, name string) string {
	start := time.Now()
	greeting := l.next.Greet(accountID, name)

	l.logger.Debug().
		Str("method", "Hello.Greet").
		Dur("duration", time.Since(start)).
		Str("greeting", greeting).
		Msg("")

	return greeting
}

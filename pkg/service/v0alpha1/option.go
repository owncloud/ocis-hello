package v0alpha1svc

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Option configures an endpoint option.
type Option func(*options)

// WithLogger returns an option to set the logger.
func WithLogger(val log.Logger) Option {
	return func(o *options) {
		o.logger = val
	}
}

// WithCounter returns an option to set the counter.
func WithCounter(val metrics.Counter) Option {
	return func(o *options) {
		o.counter = val
	}
}

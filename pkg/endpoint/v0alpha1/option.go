package v0alpha1endpoint

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/owncloud/ocis-hello/pkg/service/v0alpha1"
)

// Option configures an endpoint option.
type Option func(*options)

// WithService returns an option to set the service.
func WithService(val v0alpha1svc.Service) Option {
	return func(o *options) {
		o.service = val
	}
}

// WithLogger returns an option to set the logger.
func WithLogger(val log.Logger) Option {
	return func(o *options) {
		o.logger = val
	}
}

// WithDuration returns an option to set the duration.
func WithDuration(val metrics.Histogram) Option {
	return func(o *options) {
		o.duration = val
	}
}

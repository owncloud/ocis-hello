package assets

import (
	"github.com/go-kit/kit/log"
)

// Option configures an assets option.
type Option func(*assets)

// WithLogger returns an option to set the logger.
func WithLogger(val log.Logger) Option {
	return func(a *assets) {
		a.logger = val
	}
}

// WithPath returns an option to set the path.
func WithPath(val string) Option {
	return func(a *assets) {
		a.path = val
	}
}

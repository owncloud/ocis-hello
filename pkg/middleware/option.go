package middleware

import (
	"github.com/go-kit/kit/log"
)

type options struct {
	logger log.Logger
	root   string
	assets string
}

// Option configures options.
type Option func(*options)

// WithLogger returns an option to set logger.
func WithLogger(val log.Logger) Option {
	return func(s *options) {
		s.logger = val
	}
}

// WithRoot returns an option to set root.
func WithRoot(val string) Option {
	return func(s *options) {
		s.root = val
	}
}

// WithAssets returns an option to set assets.
func WithAssets(val string) Option {
	return func(s *options) {
		s.assets = val
	}
}

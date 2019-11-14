package debug

import (
	"github.com/go-kit/kit/log"
)

// Option configures an Server option.
type Option func(*Server)

// WithLogger returns an option to set logger.
func WithLogger(val log.Logger) Option {
	return func(s *Server) {
		s.logger = val
	}
}

// WithAddr returns an option to set addr.
func WithAddr(val string) Option {
	return func(s *Server) {
		s.addr = val
	}
}

// WithToken returns an option to set token.
func WithToken(val string) Option {
	return func(s *Server) {
		s.token = val
	}
}

// WithPprof returns an option to set pprof.
func WithPprof(val bool) Option {
	return func(s *Server) {
		s.pprof = val
	}
}

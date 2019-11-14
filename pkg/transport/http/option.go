package http

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

// WithGrpc returns an option to set addr.
func WithGrpc(val string) Option {
	return func(s *Server) {
		s.grpc = val
	}
}

// WithRoot returns an option to set root.
func WithRoot(val string) Option {
	return func(s *Server) {
		s.root = val
	}
}

// WithAssets returns an option to set assets.
func WithAssets(val string) Option {
	return func(s *Server) {
		s.assets = val
	}
}

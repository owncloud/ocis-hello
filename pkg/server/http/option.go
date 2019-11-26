package http

import (
	"github.com/rs/zerolog"
)

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type Option func(o *Options)

type Options struct {
	Logger zerolog.Logger
	Addr   string
	Assets string
}

func Logger(val zerolog.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

func Addr(val string) Option {
	return func(o *Options) {
		o.Addr = val
	}
}

func Assets(val string) Option {
	return func(o *Options) {
		o.Assets = val
	}
}

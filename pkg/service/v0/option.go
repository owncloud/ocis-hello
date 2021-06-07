package service

import (
	"github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger       log.Logger
	Config       *config.Config
	PhraseSource GreetingPhraseSource
}

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the Logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// PhraseSource provides the phrase source for the greeter service.
func PhraseSource(src GreetingPhraseSource) Option {
	return func(o *Options) {
		o.PhraseSource = src
	}
}

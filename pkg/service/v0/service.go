package service

import (
	"errors"
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

const (
	// DefaultPhrase defines the default phrase
	DefaultPhrase = "Hello %s"
)

var (
	ErrMissingName = errors.New("name missing")
)

type Greeter interface {
	Greet(accountID, name string) (greeting string)
}

type GreetingPhraseSource interface {
	GetPhrase(accountID string) (phrase string)
}

type StaticPhraseSource struct {
	Phrase string
}

func (s StaticPhraseSource) GetPhrase(accountID string) string {
	return s.Phrase
}

// New returns a new instance of Service
func NewGreeter(opts ...Option) (Greeter, error) {
	options := newOptions(opts...)

	g := BasicGreeter{
		log:          options.Logger,
		phraseSource: options.PhraseSource,
	}

	return g, nil
}

// BasicGreeter implements the Greeter interface
type BasicGreeter struct {
	log          log.Logger
	phraseSource GreetingPhraseSource
}

// Greet implements the HelloHandler interface.
func (g BasicGreeter) Greet(accountID, name string) string {
	phrase := g.phraseSource.GetPhrase(accountID)
	return fmt.Sprintf(phrase, name)
}

package v0alpha1svc

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/metrics"
)

// Instrument returns a service middleware that instruments service metrics.
func Instrument(counter metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrument{
			next:    next,
			counter: counter,
		}
	}
}

type instrument struct {
	next    Service
	counter metrics.Counter
}

func (i instrument) Greet(ctx context.Context, name string) (string, error) {
	res, err := i.next.Greet(ctx, name)

	i.counter.With(
		"method",
		"greet",
	).With(
		"success",
		fmt.Sprint(err == nil),
	).Add(
		1.0,
	)

	return res, err
}

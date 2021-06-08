package service

import (
	"context"

	"go.opencensus.io/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Greeter) Greeter {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Greeter
}

// Greet implements the Greeter interface.
func (t tracing) Greet(accountID, name string) string {
	_, span := trace.StartSpan(context.Background(), "Hello.Greet")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("name", name),
		trace.StringAttribute("accountID", accountID),
	}, "Execute Hello.Greet handler")

	return t.next.Greet(accountID, name)
}

package svc

import (
	"context"

	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"go.opencensus.io/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next v0proto.HelloHandler) v0proto.HelloHandler {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next v0proto.HelloHandler
}

// Greet implements the HelloHandler interface.
func (t tracing) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	ctx, span := trace.StartSpan(ctx, "Hello.Greet")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("name", req.Name),
	}, "Execute Hello.Greet handler")

	return t.next.Greet(ctx, req, rsp)
}

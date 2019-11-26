package svc

import (
	"context"
	"errors"
	"fmt"

	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
)

var (
	ErrMissingName = errors.New("missing a name")
)

type Hello struct {
}

func (s *Hello) Greet(ctx context.Context, req *v0proto.GreetRequest, rsp *v0proto.GreetResponse) error {
	if req.Name == "" {
		return ErrMissingName
	}

	rsp.Message = fmt.Sprintf(
		"Hello %s",
		req.Name,
	)

	return nil
}

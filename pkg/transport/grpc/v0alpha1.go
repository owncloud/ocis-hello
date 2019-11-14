package grpc

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/owncloud/ocis-hello/pkg/api/v0alpha1"
	"github.com/owncloud/ocis-hello/pkg/endpoint/v0alpha1"
)

type v0alpha1 struct {
	greet grpc.Handler
}

func NewV0alpha1(endpoints v0alpha1endpoint.Set, logger log.Logger) v0alpha1api.HelloServiceServer {
	return &v0alpha1{
		greet: grpc.NewServer(
			endpoints.GreetEndpoint,
			decodeGRPCGreetRequest,
			encodeGRPCGreetResponse,
			grpc.ServerErrorHandler(
				transport.NewLogErrorHandler(logger),
			),
		),
	}
}

func (s *v0alpha1) Greet(ctx context.Context, req *v0alpha1api.GreetRequest) (*v0alpha1api.GreetResponse, error) {
	_, rep, err := s.greet.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*v0alpha1api.GreetResponse), nil
}

func decodeGRPCGreetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*v0alpha1api.GreetRequest)

	return v0alpha1endpoint.GreetRequest{
		Name: req.Name,
	}, nil
}

func decodeGRPCGreetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*v0alpha1api.GreetResponse)

	return v0alpha1endpoint.GreetResponse{
		Message: reply.Message,
		Err:     str2err(reply.Err),
	}, nil
}

func encodeGRPCGreetResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(v0alpha1endpoint.GreetResponse)

	return &v0alpha1api.GreetResponse{
		Message: resp.Message,
		Err:     err2str(resp.Err),
	}, nil
}

func encodeGRPCGreetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(v0alpha1endpoint.GreetRequest)

	return &v0alpha1api.GreetRequest{
		Name: req.Name,
	}, nil
}
